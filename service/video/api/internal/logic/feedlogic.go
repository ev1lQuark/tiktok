package logic

import (
	"context"
	"fmt"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"strconv"
	"time"

	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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
	tableVideos, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.PublishTime.Gt(lastTime)).Order(videoQuery.PublishTime.Desc()).Limit(l.svcCtx.Config.Video.NumberLimit).Find()

	if err != nil {
		msg := fmt.Sprintf("查询视频失败：%v", err)
		logx.Error(msg)
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}
	authorIds := make([]int64, 0, len(tableVideos))
	videoIds := make([]int64, 0, len(tableVideos))

	for _, value := range tableVideos {
		authorIds = append(authorIds, value.AuthorID)
		videoIds = append(videoIds, value.ID)
	}

	var eg errgroup.Group

	//根据userId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetTotalFavoriteNumReply

	eg.Go(func() error {
		var err error
		totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetTotalFavoriteNum(l.ctx, &like.GetTotalFavoriteNumReq{UserId: authorIds})
		return err
	})

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
	eg.Go(func() error {
		var err error
		userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(l.ctx, &like.GetFavoriteCountByUserIdReq{UserId: authorIds})
		return err
	})

	// 根据videoId获取视屏点赞总数
	var videoFavoriteCountList *like.GetFavoriteCountByVideoIdReply
	eg.Go(func() error {
		var err error
		videoFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByVideoId(l.ctx, &like.GetFavoriteCountByVideoIdReq{VideoId: videoIds})
		return err
	})

	// 根据videoId获取视屏评论总数
	var videoCommentCountList *comment.GetComentCountByVideoIdReply
	eg.Go(func() error {
		var err error
		videoCommentCountList, err = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{VideoId: videoIds})
		return err
	})

	// 根据userId和videoId判断是否点赞
	var isFavoriteList = &like.IsFavoriteReply{}
	eg.Go(func() error {
		var err error
		isFavoriteList, err = l.svcCtx.LikeRpc.IsFavorite(l.ctx, &like.IsFavoriteReq{VideoId: videoIds, UserId: authorIds})
		return err
	})

	//错误判断
	if err := eg.Wait(); err != nil {
		msg := fmt.Sprintf("調用Rpc失敗%v", err)
		logx.Error(msg)
		resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
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
	fmt.Println(len(userNameList.NameList))
	fmt.Println(len(totalFavoriteNumList.Count))
	fmt.Println(len(userFavoriteCountList.Count))
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
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
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

	fmt.Println(99999)
	nextTime := 0
	if len(videos) != 0 {
		nextTime = int(tableVideos[len(tableVideos)-1].PublishTime.Unix())
	}
	resp = &types.FeedReply{StatusCode: res.SuccessCode, StatusMsg: "请求成功", NextTime: nextTime, VideoList: videos}
	return resp, nil
}
