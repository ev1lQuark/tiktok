package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"path"
	"strconv"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"golang.org/x/sync/errgroup"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListReq) (resp *types.PublishListReply, err error) {
	// 登录校验
	if ok := jwt.Verify(l.svcCtx.Config.Auth.AccessSecret, req.Token); !ok {
		resp = &types.PublishListReply{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}
	// 参数校验
	if len(req.UserID) == 0 {
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	userId, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	// 查找last date最近视频
	videoQuery := l.svcCtx.Query.Video

	tableVideos, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.AuthorID.Eq(userId)).Order(videoQuery.PublishTime.Desc()).Find()

	if err != nil {
		log.Printf("获取用户的视频发布列表失败：%v", err)
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: "获取用户的视频发布列表失败"}
		return resp, nil
	}

	videoIds := make([]int64, 0, len(tableVideos))
	for _, value := range tableVideos {
		videoIds = append(videoIds, value.ID)
	}
	authorIds := []int64{userId}

	var eg errgroup.Group

	//根据userId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetFavoriteCountByAuthorIdsReply
	eg.Go(func() error {
		var err error
		totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetFavoriteCountByAuthorIds(l.ctx, &like.GetFavoriteCountByAuthorIdsReq{AuthorIds: authorIds})
		return err
	})

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdsReply
	eg.Go(func() error {
		var err error
		userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserIds(l.ctx, &like.GetFavoriteCountByUserIdsReq{UserIds: authorIds})
		return err
	})

	// 根据videoId获取视频点赞总数
	var videoFavoriteCountList *like.GetFavoriteCountByVideoIdsReply
	eg.Go(func() error {
		var err error
		videoFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByVideoIds(l.ctx, &like.GetFavoriteCountByVideoIdsReq{VideoIds: videoIds})
		return err
	})

	// 根据videoId获取视频评论总数
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
		authorIds := make([]int64, 0, len(tableVideos))
		for range tableVideos {
			authorIds = append(authorIds, userId)
		}
		isFavoriteList, err = l.svcCtx.LikeRpc.IsFavorite(l.ctx, &like.IsFavoriteReq{VideoIds: videoIds, UserIds: authorIds})
		return err
	})

	// errorgroup 等待所有请求完成 错误处理
	if err := eg.Wait(); err != nil {
		msg := fmt.Sprintf("調用Rpc失敗%v", err)
		logx.Error(msg)
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}

	// 获取workCount


	count, err := l.svcCtx.Redis.HGet(context.TODO(), AuthorIdToWorkCount, strconv.FormatInt(userId, 10)).Result()
	if err == redis.Nil {
		count = "0"
	} else if err != nil {
		msg := fmt.Sprintf("Redis查询失败：%v", err)
		logx.Error(msg)
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}

	workCount, _ := strconv.ParseInt(count, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("count解析int失败：%v", err)
		logx.Error(msg)
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}


	// 拼接请求
	videos := make([]types.VideoList, 0, len(tableVideos))
	for index, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			Author: types.Author{
				ID:              int(userId),
				Name:            userNameList.NameList[0],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.FormatInt(totalFavoriteNumList.CountSlice[0], 10),
				WorkCount:       int(workCount),
				FavoriteCount:   int(userFavoriteCountList.CountSlice[0]),
			},
			PlayURL:       "http://" + path.Join(l.svcCtx.Config.Minio.Endpoint, value.PlayURL),
			CoverURL:      "http://" + path.Join(l.svcCtx.Config.Minio.Endpoint, value.CoverURL),
			FavoriteCount: int(videoFavoriteCountList.CountSlice[0]),
			CommentCount:  int(videoCommentCountList.Count[index]),
			IsFavorite:    isFavoriteList.IsFavoriteSlice[index],
			Title:         value.Title,
		})
	}
	resp = &types.PublishListReply{StatusCode: res.SuccessCode, StatusMsg: "成功", VideoList: videos}
	return resp, nil
}
