package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByVideoIdLogic {
	return &GetFavoriteCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据videoId获取视屏点赞总数
func (l *GetFavoriteCountByVideoIdLogic) GetFavoriteCountByVideoId(in *like.GetFavoriteCountByVideoIdReq) (*like.EtFavoriteCountByUserIdReply, error) {
	// todo: add your logic here and delete this line

	return &like.EtFavoriteCountByUserIdReply{}, nil
}
