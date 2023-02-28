package svc

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/commentclient"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/video/query"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var VideoDataListJSON = "VIDEO::ZSET::VIDEO_DATA_JSON"
var AuthorIdToWorkCount = "VIDEO::AUTHORID::VIDEO_COUNT"

type ServiceContext struct {
	Config        config.Config
	Query         *query.Query
	MinioClient   *minio.Client
	UserRpc       userclient.User
	CommentRpc    commentclient.Comment
	LikeRpc       likeclient.Like
	Redis         *redis.Client
	ContinuedTime int64
}

func NewServiceContext(c config.Config) *ServiceContext {

	mc, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}

	svcCtx := &ServiceContext{
		Config:        c,
		Query:         query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		MinioClient:   mc,
		UserRpc:       userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CommentRpc:    commentclient.NewComment(zrpc.MustNewClient(c.CommentRpc)),
		LikeRpc:       likeclient.NewLike(zrpc.MustNewClient(c.LikeRpc)),
		Redis:         redis.NewClient(&redis.Options{Addr: c.Redis.Addr, DB: c.Redis.DB}),
		ContinuedTime: c.ContinuedTime,
	}

	err = readAllVideoCountByAuthorId(svcCtx)
	if err != nil {
		return nil
	}

	go TimedTask(svcCtx)

	return svcCtx
}

func TimedTask(svcCtx *ServiceContext) error {
	// 定时删除冷数据
	for {
		maxTime := time.Now().Add(-time.Hour * time.Duration(svcCtx.ContinuedTime))
		_, err := svcCtx.Redis.ZRemRangeByScore(context.TODO(), VideoDataListJSON, "0", strconv.FormatInt(maxTime.Unix(), 10)).Result()
		if err != nil {
			logx.Error("定时删除冷数据失败%w", err)
		}
		time.Sleep(time.Hour * time.Duration(svcCtx.ContinuedTime))
	}

}

func readAllVideoCountByAuthorId(svcCtx *ServiceContext) error {
	videoQuery := svcCtx.Query.Video
	videoList, err := videoQuery.WithContext(context.TODO()).Select(videoQuery.AuthorID).Find()
	if err != nil {
		return err
	}

	videoCount := make(map[int64]int64)
	for _, video := range videoList {
		videoCount[video.AuthorID]++
	}

	for key, value := range videoCount {
		_, err = svcCtx.Redis.HSet(context.TODO(), AuthorIdToWorkCount, strconv.FormatInt(key, 10), strconv.FormatInt(value, 10)).Result()
		if err != nil {
			return err
		}
	}

	return nil

}
