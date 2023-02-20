package logic

import (
	"context"
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
func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *comment.GetComentCountByVideoIdReq) (*comment.GetComentCountByVideoIdReply, error) {
	// todo: add your logic here and delete this line

	commentQuery := l.svcCtx.Query.Comment
	numList := make([]int64, 0, len(in.VideoId))
	for _, videoId := range in.VideoId {
		num, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Count()
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &comment.GetComentCountByVideoIdReply{Count: numList}, nil
}
