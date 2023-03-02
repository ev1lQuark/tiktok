package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/like/pattern"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeListLogic {
	return &GetLikeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 根据userId获取点赞video列表
func (l *GetLikeListLogic) GetLikeList(req *types.LikeListRequest) (resp *types.LikeListResponse, err error) {
	// Parse jwt token
	if ok := jwt.Verify(l.svcCtx.Config.Auth.AccessSecret, req.Token); !ok {
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.AuthFailedCode),
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}
	// 参数校验
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		logx.Errorf("param error: %w", err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.AuthFailedCode),
			StatusMsg:  "参数错误",
		}
		return resp, nil
	}

	// 获取 userId 喜欢视频的 videoIds
	likeVideoIds, err := l.svcCtx.Redis.SMembers(l.ctx, pattern.GetLikeSetUserIdKey(userId)).Result()
	if err != nil {
		logx.Errorf("redis error: %w", err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.AuthFailedCode),
			StatusMsg:  "redis 错误",
		}
		return resp, err
	}
	if len(likeVideoIds) == 0 {
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.SuccessCode),
			StatusMsg:  "success",
			VideoList:  make([]types.VideoList, 0),
		}
		return resp, nil
	}

	// 格式转换
	videoIds := make([]int64, 0, len(likeVideoIds))
	for i := range likeVideoIds {
		videoId, _ := strconv.ParseInt(likeVideoIds[i], 10, 64)
		videoIds = append(videoIds, videoId)
	}

	// 获取视频Infos
	videosInfoReply, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{
		VideoId: videoIds,
	})
	if err != nil {
		logx.Errorf("rpc error: %w", err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.RemoteServiceErrorCode),
			StatusMsg:  err.Error(),
		}
		return resp, nil
	}
	// 整理出authorIds
	authorIds := videosInfoReply.AuthorId

	var eg errgroup.Group

	// 根据 videoIds 获取 commentCount
	var commentCountReply *comment.GetComentCountByVideoIdReply
	eg.Go(func() error {
		var err error
		commentCountReply, err = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{
			VideoId: videoIds,
		})
		if err != nil {
			return fmt.Errorf("CommentRpc.GetCommentCountByVideoId error: %v", err)
		}
		return nil
	})

	// 根据 authorIds 获取 workCount
	var workCount *video.VideoNumReply
	eg.Go(func() error {
		var err error
		workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{AuthorId: authorIds})
		if err != nil {
			return fmt.Errorf("VideoRpc.GetVideoNumByAuthorId error: %v", err)
		}
		return nil
	})

	//根据 authorIds 获取 authorNames
	var authorNamesReply *user.NameListReply
	eg.Go(func() error {
		var err error
		authorNamesReply, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		if err != nil {
			return fmt.Errorf("UserRpc.GetNames error: %v", err)
		}
		return nil
	})

	// errgroup 等待所有 rpc 请求完成
	if err := eg.Wait(); err != nil {
		msg := fmt.Sprintf("rpc error: %s", err.Error())
		logx.Error(msg)
		resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.InternalServerErrorCode), StatusMsg: msg}
		return resp, err
	}

	authorFavoriteCountList := make([]int64, 0)
	authorIsFavoritedCountList := make([]int64, 0)
	pipe := l.svcCtx.Redis.Pipeline()
	for _, authorId := range authorIds {
		res, _ := pipe.HGet(l.ctx, pattern.LikeMapUserIdCountKey, strconv.FormatInt(authorId, 10)).Result()
		count, _ := strconv.ParseInt(res, 10, 64)
		authorFavoriteCountList = append(authorFavoriteCountList, count)
		res, _ = pipe.HGet(l.ctx, pattern.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10)).Result()
		count, _ = strconv.ParseInt(res, 10, 64)
		authorIsFavoritedCountList = append(authorFavoriteCountList, count)
	}
	_, err = pipe.Exec(l.ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logx.Error(err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.InternalServerErrorCode),
			StatusMsg:  err.Error(),
		}
		return resp, nil
	}

	videoList := make([]types.VideoList, 0, len(videoIds))
	pipe = l.svcCtx.Redis.Pipeline()
	for i, videoId := range videoIds {
		//通过videoId获取当前视频受喜欢次数
		res, _ := pipe.HGet(l.ctx, pattern.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10)).Result()
		videoFavoriteCount, _ := strconv.ParseInt(res, 10, 64)
		//通过videoId判断用户是否对其点赞
		isF, _ := pipe.SIsMember(l.ctx, pattern.GetLikeSetUserIdKey(userId), strconv.FormatInt(videoId, 10)).Result()
		//对每个video进行整理,
		videoSingle := types.VideoList{
			ID: videoId,
			Author: types.Author{
				ID:              authorIds[i],
				Name:            authorNamesReply.NameList[i],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.FormatInt(authorIsFavoritedCountList[i], 10),
				WorkCount:       workCount.VideoNum[i],
				FavoriteCount:   authorFavoriteCountList[i],
			},
			PlayURL:       videosInfoReply.PlayUrl[i],
			CoverURL:      videosInfoReply.CoverUrl[i],
			FavoriteCount: videoFavoriteCount,
			CommentCount:  commentCountReply.Count[i],
			IsFavorite:    isF,
			Title:         videosInfoReply.Title[i],
		}
		videoList = append(videoList, videoSingle)
	}
	resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.SuccessCode), StatusMsg: "", VideoList: videoList}
	return resp, nil
}
