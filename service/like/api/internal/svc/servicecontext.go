package svc

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/commentclient"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/like/query"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

var MsgPattern = "userId:%d,videoId:%d,authorId:%d,cancel:%d"

type ServiceContext struct {
	Config     config.Config
	Query      *query.Query
	UserRpc    userclient.User
	CommentRpc commentclient.Comment
	VideoRpc   videoclient.Video
	Redis      *redis.Client
	MqProducer rocketmq.Producer
	MqConsumer rocketmq.PushConsumer
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
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CommentRpc: commentclient.NewComment(zrpc.MustNewClient(c.CommentRpc)),
		VideoRpc:   videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		Redis:      redis.NewClient(&redis.Options{Addr: c.Redis.Addr, DB: c.Redis.DB}),
		MqProducer: pdc,
		MqConsumer: csm,
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
				var userId, videoId, authorId int64
				var cancel int32
				_, err := fmt.Sscanf(string(msg.Body), MsgPattern, &userId, &videoId, &authorId, &cancel)
				if err != nil {
					logx.Error(err)
					return consumer.ConsumeRetryLater, err
				}
				// 修改数据库
				likeQuery := svcCtx.Query.Like
				like, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.VideoID.Eq(videoId)).FirstOrCreate()
				if err != nil {
					logx.Errorf("查询数据库失败%w", err)
					return consumer.ConsumeRetryLater, err
				}

				_, err = likeQuery.WithContext(context.TODO()).Where(likeQuery.ID.Eq(like.ID)).UpdateSimple(likeQuery.Cancel.Value(cancel), likeQuery.AuthorID.Value(authorId))
				if err != nil {
					logx.Error(err)
					return consumer.ConsumeRetryLater, err
				}
			}
			return consumer.ConsumeSuccess, nil
		})
	c.Start()
}
