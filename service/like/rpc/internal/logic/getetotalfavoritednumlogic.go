package logic

import (
	"context"
	"gorm.io/gorm"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeteTotalFavoritedNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGeteTotalFavoritedNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeteTotalFavoritedNumLogic {
	return &GeteTotalFavoritedNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号所发视频获赞总数
func (l *GeteTotalFavoritedNumLogic) GeteTotalFavoritedNum(in *like.GetFavoriteCountByUserIdReq) (*like.GeteTotalFavoritedNumReply, error) {

	authorId := in.UserId
	likeQuery := l.svcCtx.Query.Like
	num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.AuthorID.Eq(authorId[0])).Count()
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
	return &like.GeteTotalFavoritedNumReply{Count: count}, nil
}
