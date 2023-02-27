package svc

import (
	"context"
	"strconv"
	"time"

	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/commentclient"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/like/query"
	"github.com/ev1lQuark/tiktok/service/like/setting"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
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
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:     c,
		Query:      query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CommentRpc: commentclient.NewComment(zrpc.MustNewClient(c.CommentRpc)),
		VideoRpc:   videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		Redis:      redis.NewClient(&redis.Options{Addr: c.Redis.Addr, DB: c.Redis.DB}),
	}

	cronInstance := cron.New()
	cronInstance.AddFunc("@daily", func() { sync2db(svcCtx) })
	cronInstance.Start()

	return svcCtx
}

func sync2db(svcCtx *ServiceContext) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	likeQuery := svcCtx.Query.Like
	m, err := svcCtx.Redis.HGetAll(ctx, setting.LikeMapDataKey).Result()
	if err != nil {
		logx.Error(err)
	}
	for k, v := range m {
		userId, videoId := setting.ParseLikeMapDataKey(k)
		cancel, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logx.Error(err)
			continue
		}
		like, err := likeQuery.WithContext(ctx).Where(likeQuery.UserID.Eq(userId), likeQuery.VideoID.Eq(videoId)).FirstOrCreate()
		if err != nil {
			logx.Error(err)
			continue
		}
		if like.Cancel != int32(cancel) {
			_, err = likeQuery.WithContext(ctx).Update(likeQuery.Cancel, cancel)
			if err != nil {
				logx.Error(err)
			}
		}
	}
}
