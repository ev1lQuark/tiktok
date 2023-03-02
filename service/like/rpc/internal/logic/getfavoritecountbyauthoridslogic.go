package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/pattern"
	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

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
	authorIds := make([]string, 0, len(in.AuthorIds))
	for _, authorId := range in.AuthorIds {
		authorIds = append(authorIds, strconv.FormatInt(authorId, 10))
	}
	res, err := l.svcCtx.Redis.HMGet(l.ctx, pattern.LikeMapAuthorIdCountKey, authorIds...).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &like.GetFavoriteCountByAuthorIdsReply{}, nil
	}
	counts := make([]int64, 0, len(res))
	for _, item := range res {
		count, _ := strconv.ParseInt(item.(string), 10, 64)
		counts = append(counts, count)
	}

	return &like.GetFavoriteCountByAuthorIdsReply{CountSlice: counts}, nil
}
