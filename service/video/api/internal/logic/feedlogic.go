package logic

import (
	"context"
	"fmt"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"strconv"
	"sync"
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

	authorIds := make([]int64, 0, len(tableVideos))
	videoIds := make([]int64, 0, len(tableVideos))
	allErrors := make([]error, 0, 10)
	for _, value := range tableVideos {
		authorIds = append(authorIds, value.AuthorID)
		videoIds = append(videoIds, value.ID)
	}

	wg := sync.WaitGroup{}
	wg.Add(6)

	// 根据userId获取userName
	var userNameList *user.NameListReply
	go func() {
		userNameList, allErrors[0] = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		wg.Done()
	}()

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetTotalFavoriteNumReply
	go func() {
		totalFavoriteNumList, allErrors[1] = l.svcCtx.LikeRpc.GetTotalFavoriteNum(l.ctx, &like.GetTotalFavoriteNumReq{UserId: authorIds})
		wg.Done()
	}()

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
	go func() {
		userFavoriteCountList, allErrors[2] = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(l.ctx, &like.GetFavoriteCountByUserIdReq{UserId: authorIds})
		wg.Done()
	}()

	// 根据videoId获取视屏点赞总数
	var videoFavoriteCountList *like.GetFavoriteCountByVideoIdReply
	go func() {
		videoFavoriteCountList, allErrors[3] = l.svcCtx.LikeRpc.GetFavoriteCountByVideoId(l.ctx, &like.GetFavoriteCountByVideoIdReq{VideoId: videoIds})
		wg.Done()
	}()

	// 根据videoId获取视屏评论总数
	var videoCommentCountList *comment.GetComentCountByVideoIdReply
	go func() {
		videoCommentCountList, allErrors[4] = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{VideoId: videoIds})
		wg.Done()
	}()

	// 根据userId和videoId判断是否点赞
	var isFavoriteList *like.IsFavoriteReply
	go func() {
		isFavoriteList, allErrors[5] = l.svcCtx.LikeRpc.IsFavorite(l.ctx, &like.IsFavoriteReq{VideoId: videoIds, UserId: authorIds})
		wg.Done()
	}()

	wg.Wait()

	//错误判断
	for _, err := range allErrors {
		if err != nil {
			msg := fmt.Sprintf("調用Rpc失敗%v", err)
			logx.Error(msg)
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
			return resp, nil
		}
	}

	// 获取workCount
	workCount := make([]int, 0, len(tableVideos))
	for index := 0; index < len(tableVideos); index++ {
		count, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.AuthorID.Eq(authorIds[index])).Count()
		if err != nil {
			msg := fmt.Sprintf("查询视频失败：%v", err)
			logx.Error(msg)
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
			return resp, nil
		}
		workCount = append(workCount, int(count))
	}

	// 拼接请求
	videos := make([]types.VideoList, 0, len(tableVideos))
	for index, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			Author: types.Author{
				ID:              int(authorIds[index]),
				Name:            userNameList.NameList[index],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "data:image/png;base64,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",
				BackgroundImage: "data:image/png;base64,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",
				Signature:       "愛抖音，爱生活",
				TotalFavorited:  strconv.Itoa(int(totalFavoriteNumList.Count[index])),
				WorkCount:       workCount[index],
				FavoriteCount:   int(userFavoriteCountList.Count[index]),
			},
			PlayURL:       l.svcCtx.Config.Minio.Endpoint + "/" + value.PlayURL,
			CoverURL:      l.svcCtx.Config.Minio.Endpoint + "/" + value.CoverURL,
			FavoriteCount: int(videoFavoriteCountList.Count[index]),
			CommentCount:  int(videoCommentCountList.Count[index]),
			IsFavorite:    isFavoriteList.IsFavorite[index],
			Title:         value.Title,
		})
	}
	nextTime := 0
	if len(videos) != 0 {
		nextTime = int(tableVideos[len(tableVideos)-1].PublishTime.Unix())
	}
	resp = &types.FeedReply{StatusCode: res.SuccessCode, StatusMsg: "请求成功", NextTime: nextTime, VideoList: videos}
	return resp, nil
}
