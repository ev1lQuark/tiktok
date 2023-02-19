package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

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
	// todo: add your logic here and delete this line

	return &like.IsFavoriteReply{}, nil
}
