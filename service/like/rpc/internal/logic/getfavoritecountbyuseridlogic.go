package logic

import (
	"context"
	"gorm.io/gorm"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdLogic {
	return &GetFavoriteCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号喜欢（点赞）总数
func (l *GetFavoriteCountByUserIdLogic) GetFavoriteCountByUserId(in *like.GetFavoriteCountByUserIdReq) (*like.GetFavoriteCountByUserIdReply, error) {

	userId := in.UserId
	likeQuery := l.svcCtx.Query.Like
	num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId[0])).Count()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			break
		default:
			return nil, err
		}
	}
	var count []int64
	count = append(count, num)
	return &like.GetFavoriteCountByUserIdReply{Count: count}, nil
}
