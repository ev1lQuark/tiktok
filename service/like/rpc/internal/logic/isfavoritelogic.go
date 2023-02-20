package logic

import (
	"context"
	"gorm.io/gorm"

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
	videoId := in.VideoId
	userId := in.UserId
	likeQuery := l.svcCtx.Query.Like
	isF := false
	count, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(videoId[0])).Where(likeQuery.UserID.Eq(userId[0])).Count()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			break
		default:
			return nil, err
		}
	}
	if count > 0 {
		isF = true
	}
	var isFavorite []bool
	isFavorite = append(isFavorite, isF)
	return &like.IsFavoriteReply{IsFavorite: isFavorite}, nil
}
