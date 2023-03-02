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

type GetFavoriteCountByAuthorIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByAuthorIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByAuthorIdsLogic {
	return &GetFavoriteCountByAuthorIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号所发视频获赞总数
func (l *GetFavoriteCountByAuthorIdsLogic) GetFavoriteCountByAuthorIds(in *like.GetFavoriteCountByAuthorIdsReq) (*like.GetFavoriteCountByAuthorIdsReply, error) {
	counts := make([]int64, 0, len(in.AuthorIds))
	for _, authorId := range in.AuthorIds {
		res, err := l.svcCtx.Redis.HGet(l.ctx, pattern.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10)).Result()
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

	return &like.GetFavoriteCountByAuthorIdsReply{CountSlice: counts}, nil
}
