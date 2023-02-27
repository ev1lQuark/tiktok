package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/like/setting"

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
	videoIds := make([]string, 0, len(in.VideoIds))
	for _, videoId := range in.VideoIds {
		videoIds = append(videoIds, strconv.FormatInt(videoId, 10))
	}

	userIds := make([]string, 0, len(in.UserIds))
	for _, userId := range in.UserIds {
		userIds = append(userIds, strconv.FormatInt(userId, 10))
	}

	keys := make([]string, 0, len(videoIds))
	for i := range videoIds {
		uid, _ := strconv.ParseInt(userIds[i], 10, 64)
		vid, _ := strconv.ParseInt(videoIds[i], 10, 64)
		keys = append(keys, setting.GetLikeMapDataKey(uid, vid))
	}

	res, err := l.svcCtx.Redis.HMGet(l.ctx, setting.LikeMapDataKey, keys...).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &like.IsFavoriteReply{}, nil
	}

	isFavoriteSlice := make([]bool, 0, len(res))
	for _, item := range res {
		if item == 1 {
			isFavoriteSlice = append(isFavoriteSlice, false)
		} else {
			isFavoriteSlice = append(isFavoriteSlice, true)
		}
	}
	return &like.IsFavoriteReply{IsFavoriteSlice: isFavoriteSlice}, nil
}
