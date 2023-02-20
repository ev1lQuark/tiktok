package logic

import (
	"context"
	"gorm.io/gorm"

	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据videoId获取视屏评论总数
func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *comment.GetFavoriteCountByVideoIdReq) (*comment.GetFavoriteCountByVideoIdReply, error) {
	// todo: add your logic here and delete this line

	videoId := in.VideoId
	commentQuery := l.svcCtx.Query.Comment
	num, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId[0])).Count()
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
	return &comment.GetFavoriteCountByVideoIdReply{Count: count}, nil
}
