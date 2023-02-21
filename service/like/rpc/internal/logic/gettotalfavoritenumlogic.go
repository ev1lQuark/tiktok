package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTotalFavoriteNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTotalFavoriteNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTotalFavoriteNumLogic {
	return &GetTotalFavoriteNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号所发视频获赞总数
func (l *GetTotalFavoriteNumLogic) GetTotalFavoriteNum(in *like.GetTotalFavoriteNumReq) (*like.GetTotalFavoriteNumReply, error) {
	// todo: add your logic here and delete this line

	likeQuery := l.svcCtx.Query.Like
	numList := make([]int64, 0, len(in.UserId))
	for _, userId := range in.UserId {
		num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.AuthorID.Eq(userId)).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &like.GetTotalFavoriteNumReply{Count: numList}, nil
}
