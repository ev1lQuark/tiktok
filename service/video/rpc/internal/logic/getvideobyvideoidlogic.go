package logic

import (
	"context"
	"log"

	"github.com/ev1lQuark/tiktok/service/video/model"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByVideoIdLogic {
	return &GetVideoByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByVideoIdLogic) GetVideoByVideoId(in *video.VideoIdReq) (*video.VideoInfoReply, error) {
	videoQuery := l.svcCtx.Query.Video
	authorIdList := make([]int64, 0, len(in.VideoId))
	playUrlList := make([]string, 0, len(in.VideoId))
	coverUrlList := make([]string, 0, len(in.VideoId))
	publishTimeList := make([]string, 0, len(in.VideoId))
	tileList := make([]string, 0, len(in.VideoId))

	for _, videoId := range in.VideoId {
		video, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.ID.Eq(videoId)).Find()
		if err != nil {
			log.Printf("Rpc获取视频详信息失败：%v", err)
			return nil, err
		}
		if len(video) == 0 {
			video = append(video, &model.Video{})
		}
		authorIdList = append(authorIdList, video[0].AuthorID)
		playUrlList = append(playUrlList, video[0].PlayURL)
		coverUrlList = append(coverUrlList, video[0].CoverURL)
		publishTimeList = append(publishTimeList, video[0].PublishTime.Format("2006-01-02 15:04:05"))
		tileList = append(tileList, video[0].Title)
	}

	return &video.VideoInfoReply{
		AuthorId:    authorIdList,
		PlayUrl:     playUrlList,
		CoverUrl:    coverUrlList,
		PublishTime: publishTimeList,
		Title:       tileList,
	}, nil
}
