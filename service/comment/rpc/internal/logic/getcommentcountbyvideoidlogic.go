package logic

import (
	"context"
	"errors"

	"github.com/ev1lQuark/tiktok/service/comment/pattern"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"

	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据videoId获取视屏评论总数
func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *comment.GetComentCountByVideoIdReq) (*comment.GetComentCountByVideoIdReply, error) {
	numList := make([]int64, 0, len(in.VideoId))
	for _, videoId := range in.VideoId {
		num, err := l.svcCtx.Redis.HGet(context.TODO(), pattern.VideoIDToCommentCount, strconv.FormatInt(videoId, 10)).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				numList = append(numList, 0)
				continue
			}
			return nil, err
		}
		count, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, err
		}
		numList = append(numList, count)
	}
	return &comment.GetComentCountByVideoIdReply{Count: numList}, nil
}
