package logic

import (
	"context"
	"errors"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/pattern"
	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/redis/go-redis/v9"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByVideoIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByVideoIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByVideoIdsLogic {
	return &GetFavoriteCountByVideoIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据videoId获取视频点赞总数
func (l *GetFavoriteCountByVideoIdsLogic) GetFavoriteCountByVideoIds(in *like.GetFavoriteCountByVideoIdsReq) (*like.GetFavoriteCountByVideoIdsReply, error) {
	counts := make([]int64, 0, len(in.VideoIds))
	for _, videoId := range in.VideoIds {
		res, err := l.svcCtx.Redis.HGet(l.ctx, pattern.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10)).Result()
		var count int64
		if err != nil {
			if errors.Is(err, redis.Nil) {
				count = 0
			} else {
				return nil, err
			}
		} else {
			count, _ = strconv.ParseInt(res, 10, 64)
		}
		counts = append(counts, count)
	}

	return &like.GetFavoriteCountByVideoIdsReply{CountSlice: counts}, nil
}
