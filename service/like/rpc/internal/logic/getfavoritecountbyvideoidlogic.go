package logic

import (
	"context"
	"gorm.io/gorm"

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
	videoId := in.VideoId
	likeQuery := l.svcCtx.Query.Like
	num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(videoId[0])).Count()
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
	return &like.GetFavoriteCountByVideoIdReply{Count: count}, nil
}
