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
func (l *GetFavoriteCountByVideoIdLogic) GetFavoriteCountByVideoId(in *like.GetFavoriteCountByVideoIdReq) (*like.GetFavoriteCountByVideoIdReply, error) {
	// todo: add your logic here and delete this line

	likeQuery := l.svcCtx.Query.Like

	numList := make([]int64, 0, len(in.VideoId))
	for _, videoId := range in.VideoId {
		num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(videoId)).Count()
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &like.GetFavoriteCountByVideoIdReply{Count: numList}, nil
}
