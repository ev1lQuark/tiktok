package logic

import (
	"context"
	"errors"

	"github.com/ev1lQuark/tiktok/service/like/pattern"
	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/redis/go-redis/v9"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteLogic {
	return &IsFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId和videoId判断是否点赞
func (l *IsFavoriteLogic) IsFavorite(in *like.IsFavoriteReq) (*like.IsFavoriteReply, error) {
	isFavoriteSlice := make([]bool, 0, len(in.VideoIds))
	for i := range in.VideoIds {
		key := pattern.GetLikeMapDataKey(in.UserIds[i], in.VideoIds[i])
		res, err := l.svcCtx.Redis.HGet(l.ctx, pattern.LikeMapDataKey, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				isFavoriteSlice = append(isFavoriteSlice, false)
			} else {
				return nil, err
			}
		} else {
			if res == "0" {
				isFavoriteSlice = append(isFavoriteSlice, true)
			} else {
				isFavoriteSlice = append(isFavoriteSlice, false)
			}
		}
	}

	return &like.IsFavoriteReply{IsFavoriteSlice: isFavoriteSlice}, nil
}
