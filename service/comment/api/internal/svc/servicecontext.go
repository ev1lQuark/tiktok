package svc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/comment/pattern"
	"github.com/ev1lQuark/tiktok/service/comment/query"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	Query      *query.Query
	VideoRpc   videoclient.Video
	UserRpc    userclient.User
	LikeRpc    likeclient.Like
	Redis      *redis.Client
	MqProducer rocketmq.Producer
	MqConsumer rocketmq.PushConsumer
	Delaytime  int
}

func NewServiceContext(c config.Config) *ServiceContext {
	pdc, _ := rocketmq.NewProducer(
		producer.WithNameServer([]string{c.RocketMQ.NameServer}),
		producer.WithRetry(2),
		producer.WithGroupName(c.RocketMQ.Group),
	)
	pdc.Start()

	csm, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{c.RocketMQ.NameServer}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(c.RocketMQ.Group),
	)

	svcCtx := &ServiceContext{
		Config:     c,
		Query:      query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		VideoRpc:   videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		LikeRpc:    likeclient.NewLike(zrpc.MustNewClient(c.LikeRpc)),
		Redis:      redis.NewClient(&redis.Options{Addr: c.Redis.Addr, DB: c.Redis.DB}),
		MqProducer: pdc,
		MqConsumer: csm,
		Delaytime:  c.DelayTime,
	}
	go startMQConsumer(svcCtx)
	err := readAllCommentCountByVideoId(svcCtx)
	if err != nil {
		logx.Errorf("Redis初始化失败%w", err)
	}
	return svcCtx
}

func startMQConsumer(svcCtx *ServiceContext) {
	svcCtx.MqConsumer.Subscribe(svcCtx.Config.RocketMQ.AsyncDeleteTopic, consumer.MessageSelector{},
		func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				// 处理消息
				var commentId, videoId int64
				_, err := fmt.Sscanf(string(msg.Body), "%d-%d", &commentId, &videoId)
				if err != nil {
					msg := fmt.Sprintf("parse id from msg failed: %s", err.Error())
					logx.Error(msg)
					return consumer.ConsumeResult(consumer.FailedReturn), err
				}

				// 修改数据库
				commentQuery := svcCtx.Query.Comment
				info, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.ID.Eq(commentId)).Update(commentQuery.Cancel, 1)
				if err != nil {
					msg := fmt.Sprintf("db error: %s", err.Error())
					logx.Error(msg)
					return consumer.ConsumeResult(consumer.FailedReturn), err
				}
				if info.RowsAffected != 1 {
					msg := fmt.Sprintf("error: %s", "评论不存在")
					logx.Error(msg)
					return consumer.ConsumeResult(consumer.FailedReturn), err
				}
				logx.Info("删除评论成功")

				// 延时双删
				msg := primitive.NewMessage(svcCtx.Config.RocketMQ.ClearCacheTopic, []byte(fmt.Sprintf("%d", videoId)))
				for i := 0; i < 10; i++ { // 确保发送成功
					if err := svcCtx.MqProducer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
						if err != nil {
							logx.Error(err)
						}
					}, msg); err != nil {
						logx.Error(err)
					} else {
						break
					}
				}

			}
			return consumer.ConsumeSuccess, nil
		})

	svcCtx.MqConsumer.Subscribe(svcCtx.Config.RocketMQ.ClearCacheTopic, consumer.MessageSelector{},
		func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				// 处理消息
				var videoId int64
				_, err := fmt.Sscanf(string(msg.Body), "%d", &videoId)
				if err != nil {
					logx.Error(err)
					return consumer.ConsumeResult(consumer.FailedReturn), err
				}
				if err := svcCtx.Redis.Del(context.Background(), fmt.Sprintf(pattern.VideoIDToCommentListJSON, videoId)).Err(); err != nil {
					logx.Error(err)
					return consumer.ConsumeRetryLater, err
				} else {
					logx.Info("清除缓存成功")
				}
			}
			return consumer.ConsumeSuccess, nil
		})

	svcCtx.MqConsumer.Start()
	logx.Info("RocketMQ Consumer Started")
}

func readAllCommentCountByVideoId(svcCtx *ServiceContext) error {
	commentQuery := svcCtx.Query.Comment
	commentList, err := commentQuery.WithContext(context.TODO()).Select(commentQuery.VideoID).Where(commentQuery.Cancel.Eq(0)).Find()
	if err != nil {
		return err
	}

	commentCount := make(map[int64]int64)
	for _, comment := range commentList {
		commentCount[comment.VideoID]++
	}

	for key, value := range commentCount {
		_, err := svcCtx.Redis.HSet(context.TODO(), pattern.VideoIDToCommentCount, strconv.FormatInt(key, 10), strconv.FormatInt(value, 10)).Result()
		if err != nil {
			return err
		}
	}

	return nil

}
