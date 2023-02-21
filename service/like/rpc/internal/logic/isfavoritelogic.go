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

	likeQuery := l.svcCtx.Query.Like
	isList := make([]bool, 0, len(in.VideoId))
	for index, _ := range in.UserId {
		count, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(in.VideoId[index])).Where(likeQuery.UserID.Eq(in.UserId[index])).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			return nil, err
		}
		if count > 0 {
			isList = append(isList, true)
		} else {
			isList = append(isList, false)
		}
	}
	return &like.IsFavoriteReply{IsFavorite: isList}, nil
}
