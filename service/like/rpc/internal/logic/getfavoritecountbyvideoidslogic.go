package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/pattern"
	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

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
	videoIds := make([]string, 0, len(in.VideoIds))
	for _, videoId := range in.VideoIds {
		videoIds = append(videoIds, strconv.FormatInt(videoId, 10))
	}

	res, err := l.svcCtx.Redis.HMGet(l.ctx, pattern.LikeMapVideoIdCountKey, videoIds...).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &like.GetFavoriteCountByVideoIdsReply{}, nil
	}
	counts := make([]int64, 0, len(res))
	for _, item := range res {
		count, _ := strconv.ParseInt(item.(string), 10, 64)
		counts = append(counts, count)
	}

	return &like.GetFavoriteCountByVideoIdsReply{CountSlice: counts}, nil
}
