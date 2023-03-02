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

type GetFavoriteCountByUserIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdsLogic {
	return &GetFavoriteCountByUserIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号喜欢（点赞）总数
func (l *GetFavoriteCountByUserIdsLogic) GetFavoriteCountByUserIds(in *like.GetFavoriteCountByUserIdsReq) (*like.GetFavoriteCountByUserIdsReply, error) {
	counts := make([]int64, 0, len(in.UserIds))
	for _, userId := range in.UserIds {
		res, err := l.svcCtx.Redis.HGet(l.ctx, pattern.LikeMapUserIdCountKey, strconv.FormatInt(userId, 10)).Result()
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

	return &like.GetFavoriteCountByUserIdsReply{CountSlice: counts}, nil
}
