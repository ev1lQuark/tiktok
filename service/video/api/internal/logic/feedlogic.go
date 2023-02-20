package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedReply, err error) {
	// 参数校验
	if len(req.LatestTime) == 0 {
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	lastTime := time.Now()
	t, err := strconv.ParseInt(req.LatestTime, 10, 64)
	if err != nil {
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	lastTime = time.Unix(t, 0)
	//查找last date最近视屏
	videoQuery := l.svcCtx.Query.Video

	tableVideos, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.PublishTime.Lt(lastTime)).Order(videoQuery.PublishTime.Desc()).Limit(l.svcCtx.Config.Video.NumberLimit).Find()

	if err != nil {
		msg := fmt.Sprintf("查询视频失败：%v", err)
		logx.Error(msg)
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}

	videos := make([]types.VideoList, 0, l.svcCtx.Config.Video.NumberLimit)
	for _, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			// todo: add rpc Author
			// Author:
			PlayURL:  l.svcCtx.Config.Minio.Endpoint + "/" + value.PlayURL,
			CoverURL: l.svcCtx.Config.Minio.Endpoint + "/" + value.CoverURL,
			//todo: add rpc favorite_count comment_count
			//FavoriteCount: val.
			//CommentCount:
			//IsFavorite:
			Title: value.Title,
		})
	}
	nextTime := 0
	if len(videos) != 0 {
		nextTime = int(tableVideos[len(tableVideos)-1].PublishTime.Unix())
	}
	resp = &types.FeedReply{StatusCode: res.SuccessCode, StatusMsg: "请求成功", NextTime: nextTime, VideoList: videos}
	return resp, nil
}
