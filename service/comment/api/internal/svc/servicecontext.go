package svc

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/comment/query"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"strconv"
	"time"
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
	Delaytime int
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
		Delaytime: c.DelayTime,

	}
	go startMQConsumer(svcCtx)

	return svcCtx

}

func startMQConsumer(svcCtx *ServiceContext) {
	logx.Info("启动消息Consumer")
	// 从MQ中读取数据
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{svcCtx.Config.RocketMQ.NameServer}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(svcCtx.Config.RocketMQ.Group),
	)
	if err != nil {
		logx.Error(err)
		return
	}
	c.Subscribe(svcCtx.Config.RocketMQ.Topic, consumer.MessageSelector{},
		func(ctx context.Context,
			msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				// 处理消息
				logx.Info("receive from MQ")
				var commentId, videoId int64
				_, err := fmt.Sscanf(string(msg.Body), "%d-%d", &commentId, videoId)
				if err != nil {
					msg := fmt.Sprintf("删除评论失败：%s", err.Error())
					logx.Error(msg)
					return consumer.ConsumeRetryLater, err
				}

				// 修改数据库
				commentQuery := svcCtx.Query.Comment
				info, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.ID.Eq(commentId)).Update(commentQuery.Cancel, 1)
				if err != nil {
					msg := fmt.Sprintf("删除评论失败：%s", err.Error())
					logx.Error(msg)
					return consumer.ConsumeRetryLater, err
				}
				if info.RowsAffected != 1 {
					msg := fmt.Sprintf("删除评论失败：%s", "评论不存在")
					logx.Error(msg)
					return consumer.ConsumeRetryLater, err
				}

				// 延时双删
				go func(svcCtx *ServiceContext, videoId int64) {
					time.Sleep(time.Duration(svcCtx.Delaytime) * time.Second)
					// 直接缓存失效
					svcCtx.Redis.Del(context.TODO(),strconv.FormatInt(videoId, 10))
				}(svcCtx, videoId)
			}
			return consumer.ConsumeSuccess, nil
		})
	c.Start()
}
