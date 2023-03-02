package logic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/video/pattern"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/redis/go-redis/v9"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoNumByAuthorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoNumByAuthorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoNumByAuthorIdLogic {
	return &GetVideoNumByAuthorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoNumByAuthorIdLogic) GetVideoNumByAuthorId(in *video.AuthorIdReq) (*video.VideoNumReply, error) {
	videoNumList := make([]int64, 0, len(in.AuthorId))
	for _, authorId := range in.AuthorId {
		count, err := l.svcCtx.Redis.HGet(context.TODO(), pattern.AuthorIdToWorkCount, strconv.FormatInt(authorId, 10)).Result()
		if err == redis.Nil {
			count = "0"
		} else if err != nil {
			msg := fmt.Sprintf("Redis查询失败：%v", err)
			logx.Error(msg)
			return nil, err
		}
		countInt64, _ := strconv.ParseInt(count, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("count解析int失败：%v", err)
			logx.Error(msg)
			return nil, err
		}
		videoNumList = append(videoNumList, countInt64)
	}
	return &video.VideoNumReply{VideoNum: videoNumList}, nil
}
