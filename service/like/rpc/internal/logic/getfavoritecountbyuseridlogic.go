package logic

import (
	"context"
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

	likeQuery := l.svcCtx.Query.Like
	numList := make([]int64, 0, len(in.UserId))
	for _, userId := range in.UserId {
		num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Count()
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &like.GetFavoriteCountByUserIdReply{Count: numList}, nil
}
