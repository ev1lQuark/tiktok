package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoNumByAuthorIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoNumByAuthorIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoNumByAuthorIdLogic {
	return &GetVideoNumByAuthorIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoNumByAuthorIdLogic) GetVideoNumByAuthorId(in *video.AuthorIdReq) (*video.VideoNumReply, error) {
	videoNumList := make([]int64, 0, len(in.AuthorId))
	videoQuery := l.svcCtx.Query.Video
	for _, authorId := range in.AuthorId {
		videoNum, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.AuthorID.Eq(authorId)).Count()
		if err != nil {
			logx.Errorf("rpc查询失败%w", err)
			return nil, err
		}
		videoNumList = append(videoNumList, videoNum)
	}
	return &video.VideoNumReply{VideoNum: videoNumList}, nil
}
