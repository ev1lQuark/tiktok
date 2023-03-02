package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/ev1lQuark/tiktok/service/video/model"
	"github.com/ev1lQuark/tiktok/service/video/pattern"
	"github.com/redis/go-redis/v9"

	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"

	"github.com/ev1lQuark/tiktok/common/jwt"
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
	var lastTime time.Time
	if len(req.LatestTime) == 0 {
		lastTime = time.Now()
	} else {
		t, err := strconv.ParseInt(req.LatestTime, 10, 64)
		if err != nil {
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
			return resp, nil
		}
		if t > 9999999999 {
			t = t / 1000
		}
		lastTime = time.Unix(t, 0)
	}

	var userId int64
	if req.Token != "" {
		userId, err = jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
		if err != nil {
			logx.Error(err)
			msg := fmt.Sprintf("token 验证失败：%v", err)
			return &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}, nil
		}
	}

	//查找last date最近视屏

	var tableVideos []*model.Video

	// 多设置2小时防止边界问题
	if time.Now().Sub(lastTime).Hours() > float64(l.svcCtx.Config.ContinuedTime+2) {
		videoQuery := l.svcCtx.Query.Video
		tableVideos, err = videoQuery.WithContext(context.TODO()).Where(videoQuery.PublishTime.Lt(lastTime)).Order(videoQuery.PublishTime.Desc()).Limit(l.svcCtx.Config.Video.NumberLimit).Find()
		if err != nil {
			msg := fmt.Sprintf("查询视频失败：%v", err)
			logx.Error(msg)
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
			return resp, nil
		}
	} else {
		opt := &redis.ZRangeBy{
			Min:    "0",                                      //最小分数
			Max:    strconv.FormatInt(lastTime.Unix(), 10),   //最大分数
			Offset: 0,                                        //在满足条件的范围，从offset下标处开始取值
			Count:  int64(l.svcCtx.Config.Video.NumberLimit), //查询结果集个数
		}
		tableVideoListJSON := l.svcCtx.Redis.ZRevRangeByScore(context.TODO(), pattern.VideoDataListJSON, opt).Val()
		for _, videoJSON := range tableVideoListJSON {
			var video model.Video
			err := json.Unmarshal([]byte(videoJSON), &video)
			if err != nil {
				msg := fmt.Sprintf("json反序列化失败：%v", err)
				logx.Error(msg)
				resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
				return resp, nil
			}
			tableVideos = append(tableVideos, &video)
		}
		if len(tableVideos) < l.svcCtx.Config.Video.NumberLimit {
			videoQuery := l.svcCtx.Query.Video
			tableVideos, err = videoQuery.WithContext(context.TODO()).Where(videoQuery.PublishTime.Lt(lastTime)).Order(videoQuery.PublishTime.Desc()).Limit(l.svcCtx.Config.Video.NumberLimit).Find()
			if err != nil {
				msg := fmt.Sprintf("查询视频失败：%v", err)
				logx.Error(msg)
				resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
				return resp, nil
			}
		}
	}

	authorIds := make([]int64, 0, len(tableVideos))
	videoIds := make([]int64, 0, len(tableVideos))

	userIds := make([]int64, len(tableVideos))

	for _, value := range tableVideos {
		authorIds = append(authorIds, value.AuthorID)
		videoIds = append(videoIds, value.ID)
	}

	for index := range tableVideos {
		userIds[index] = userId
	}

	var eg errgroup.Group

	//根据Id获取作者Name
	var authorNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		authorNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		if err != nil {
			return fmt.Errorf("UserRPC GetNames error: %w", err)
		}
		return nil
	})

	// 根据authorIds获取作者获赞总数
	var authorTotalFavoritedList *like.GetFavoriteCountByAuthorIdsReply
	eg.Go(func() error {
		var err error
		authorTotalFavoritedList, err = l.svcCtx.LikeRpc.GetFavoriteCountByAuthorIds(l.ctx, &like.GetFavoriteCountByAuthorIdsReq{AuthorIds: authorIds})
		if err != nil {
			return fmt.Errorf("LikeRPC GetFavoriteCountByAuthorIds error: %w", err)
		}
		return nil
	})

	// 根据authorIds 获取作者喜欢（点赞）总数
	var authorFavoriteCountList *like.GetFavoriteCountByUserIdsReply
	eg.Go(func() error {
		var err error
		authorFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserIds(l.ctx, &like.GetFavoriteCountByUserIdsReq{UserIds: authorIds})
		if err != nil {
			return fmt.Errorf("LikeRPC GetFavoriteCountByUserIds error: %w", err)
		}
		return nil
	})

	// 根据videoIds获取视频点赞总数
	var videoFavoriteCountList *like.GetFavoriteCountByVideoIdsReply
	eg.Go(func() error {
		var err error
		videoFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByVideoIds(l.ctx, &like.GetFavoriteCountByVideoIdsReq{VideoIds: videoIds})
		if err != nil {
			return fmt.Errorf("LikeRPC GetFavoriteCountByVideoIds error: %w", err)
		}
		return nil
	})

	// 根据videoId获取视频评论总数
	var videoCommentCountList *comment.GetComentCountByVideoIdReply
	eg.Go(func() error {
		var err error
		videoCommentCountList, err = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{VideoId: videoIds})
		if err != nil {
			return fmt.Errorf("CommentRPC GetCommentCountByVideoId error: %w", err)
		}
		return nil
	})

	// 根据userId和videoId判断是否点赞
	var isFavoriteList = &like.IsFavoriteReply{}
	eg.Go(func() error {
		var err error
		isFavoriteList, err = l.svcCtx.LikeRpc.IsFavorite(l.ctx, &like.IsFavoriteReq{VideoIds: videoIds, UserIds: userIds})
		if err != nil {
			return fmt.Errorf("LikeRPC IsFavorite error: %w", err)
		}
		return nil
	})

	// 错误判断
	if err := eg.Wait(); err != nil {
		logx.Error(err)
		resp = &types.FeedReply{StatusCode: res.RemoteServiceErrorCode, StatusMsg: err.Error()}
		return resp, nil
	}

	// 获取authorWorkCountList
	authorWorkCountList := make([]int, 0, len(tableVideos))
	for index := 0; index < len(tableVideos); index++ {
		count, err := l.svcCtx.Redis.HGet(context.TODO(), pattern.AuthorIdToWorkCount, strconv.FormatInt(authorIds[index], 10)).Result()

		if err == redis.Nil {
			count = "0"
		} else if err != nil {
			msg := fmt.Sprintf("Redis查询失败：%v", err)
			logx.Error(msg)
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
			return resp, nil
		}
		countInt64, _ := strconv.ParseInt(count, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("count解析int失败：%v", err)
			logx.Error(msg)
			resp = &types.FeedReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
			return resp, nil
		}
		authorWorkCountList = append(authorWorkCountList, int(countInt64))
	}

	// 拼接请求
	videos := make([]types.VideoList, 0, len(tableVideos))
	for index, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			Author: types.Author{
				ID:              int(authorIds[index]),
				Name:            authorNameList.NameList[index],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.Itoa(int(authorTotalFavoritedList.CountSlice[index])),
				WorkCount:       authorWorkCountList[index],
				FavoriteCount:   int(authorFavoriteCountList.CountSlice[index]),
			},
			PlayURL:       "http://" + path.Join(l.svcCtx.Config.Minio.Endpoint, value.PlayURL),
			CoverURL:      "http://" + path.Join(l.svcCtx.Config.Minio.Endpoint, value.CoverURL),
			FavoriteCount: int(videoFavoriteCountList.CountSlice[index]),
			CommentCount:  int(videoCommentCountList.Count[index]),
			IsFavorite:    isFavoriteList.IsFavoriteSlice[index],
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
