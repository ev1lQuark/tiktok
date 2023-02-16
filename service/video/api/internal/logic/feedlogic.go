package logic

import (
	"context"
	"fmt"
	"github.com/ev1lQuark/tiktok/common/config"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
	"time"
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

	tableVideos, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.PublishTime.Lt(lastTime)).Order(videoQuery.PublishTime.Desc()).Limit(config.VideoCount).Find()

	if err != nil {
		log.Printf("查询视屏失败：%v", err)
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: "查询视屏失败"}
		return resp, nil
	}
	log.Printf("查询成功：%v", tableVideos)

	videos := make([]types.VideoList, 0, config.VideoCount)
	for _, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			//todo add rpc Author
			//Author:
			PlayURL:  value.PlayURL,
			CoverURL: value.CoverURL,
			//todo add rpc favorite_count comment_count
			//FavoriteCount: val.
			//CommentCount:
			//IsFavorite:
			Title: value.Title,
		})
		fmt.Println(value.PublishTime)
	}
	nextTime := 0
	if len(videos) != 0 {
		nextTime = int(tableVideos[len(tableVideos)-1].PublishTime.Unix())
	}
	resp = &types.FeedReply{StatusCode: res.SuccessCode, StatusMsg: "成功", NextTime: nextTime, VideoList: videos}
	return resp, nil
}
