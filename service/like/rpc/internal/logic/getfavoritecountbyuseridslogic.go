package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/like/setting"

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
	userIds := make([]string, 0, len(in.UserIds))
	for _, userId := range in.UserIds {
		userIds = append(userIds, strconv.FormatInt(userId, 10))
	}

	res, err := l.svcCtx.Redis.HMGet(l.ctx, setting.LikeMapUserIdCountKey, userIds...).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &like.GetFavoriteCountByUserIdsReply{}, nil
	}
	counts := make([]int64, 0, len(res))
	for _, item := range res {
		count, _ := strconv.ParseInt(item.(string), 10, 64)
		counts = append(counts, count)
	}

	return &like.GetFavoriteCountByUserIdsReply{CountSlice: counts}, nil
}
