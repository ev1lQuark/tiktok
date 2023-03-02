package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/comment/model"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
)

type CommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentLogic) Comment(req *types.CommentRequest) (resp *types.CommentResponse, err error) {
	// Parse jwt token
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		logx.Error(err)
		resp = &types.CommentResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  err.Error(),
		}
		return resp, nil
	}

	// 参数校验
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误: %w", err)
		resp = &types.CommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	actionType, err := strconv.ParseInt(req.ActionType, 10, 64)
	if err != nil || (actionType != int64(1) && actionType != int64(2)) {
		logx.Errorf("参数错误: %w", err)
		resp = &types.CommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	// 判断videoId是否存在
	videoInfo, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{VideoId: []int64{videoId}})

	if err != nil {
		logx.Errorf("RPC调用失败: %w", err)
		resp = &types.CommentResponse{StatusCode: res.RemoteServiceErrorCode, StatusMsg: err.Error()}
		return resp, nil
	}
	if videoInfo.AuthorId[0] == 0 {
		logx.Errorf("invalid videoId: %d", videoId)
		resp = &types.CommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "invalid videoId"}
		return resp, nil
	}

	commentQuery := l.svcCtx.Query.Comment
	if actionType == 1 {
		// 发布评论
		comment := &model.Comment{
			UserID:      userId,
			VideoID:     videoId,
			CommentText: req.CommentText,
			CreatDate:   time.Now(),
			Cancel:      0,
		}

		// clear cache first time
		l.svcCtx.Redis.Del(l.ctx, strconv.FormatInt(videoId, 10)).Result()

		// create data in mysql
		err = commentQuery.WithContext(l.ctx).Create(comment)
		if err != nil {
			logx.Error(err)
			resp = &types.CommentResponse{StatusCode: res.InternalServerErrorCode, StatusMsg: err.Error()}
			return resp, nil
		}

		// update count
		_, err = l.svcCtx.Redis.HIncrBy(l.ctx, VideoIDToCommentCount, strconv.FormatInt(videoId, 10), 1).Result()
		if err != nil {
			logx.Error(err)
			resp = &types.CommentResponse{StatusCode: res.InternalServerErrorCode, StatusMsg: err.Error()}
			return resp, nil
		}

		// clear cache second time
		msg := primitive.NewMessage(l.svcCtx.Config.RocketMQ.ClearCacheTopic, []byte(strconv.FormatInt(videoId, 10)))
		for i := 0; i < 10; i++ { // 确保发送成功
			if err := l.svcCtx.MqProducer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
				if err != nil {
					logx.Error(err)
				}
			}, msg); err != nil {
				logx.Error(err)
			} else {
				break
			}
		}

		// RPC invoking
		var eg errgroup.Group

		// 根据userId获取userName
		var userNameList *user.NameListReply
		eg.Go(func() error {
			var err error
			userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: []int64{userId}})
			return err
		})

		// 根据userId获取本账号获赞总数
		var totalFavoriteNumList *like.GetFavoriteCountByAuthorIdsReply
		eg.Go(func() error {
			var err error
			totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetFavoriteCountByAuthorIds(l.ctx, &like.GetFavoriteCountByAuthorIdsReq{AuthorIds: []int64{userId}})
			return err
		})

		// 根据userId获取本账号喜欢（点赞）总数
		var userFavoriteCountList *like.GetFavoriteCountByUserIdsReply
		eg.Go(func() error {
			var err error
			userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserIds(l.ctx, &like.GetFavoriteCountByUserIdsReq{UserIds: []int64{userId}})
			return err
		})

		// 获取work_count
		var workCount *video.VideoNumReply
		eg.Go(func() error {
			var err error
			workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{AuthorId: []int64{userId}})
			return err
		})

		// 错误判断
		if err := eg.Wait(); err != nil {
			logx.Error(err)
			resp = &types.CommentResponse{StatusCode: res.RemoteServiceErrorCode, StatusMsg: err.Error()}
			return resp, nil
		}

		// build response
		commentInfo := types.Comment{
			ID: comment.ID,
			User: types.User{
				ID:              int(userId),
				Name:            userNameList.NameList[0],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.FormatInt(totalFavoriteNumList.CountSlice[0], 10),
				WorkCount:       int(workCount.VideoNum[0]),
				FavoriteCount:   int(userFavoriteCountList.CountSlice[0]),
			},
			Content:    comment.CommentText,
			CreateDate: comment.CreatDate.String(),
		}

		resp = &types.CommentResponse{StatusCode: res.SuccessCode, StatusMsg: "ok", Comment: commentInfo}
	} else if actionType == 2 {
		// 参数校验
		commentId, err := strconv.ParseInt(req.CommentId, 10, 64)
		if err != nil {
			logx.Error(err)
			return &types.CommentResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}, nil
		}

		// clear cache first time
		l.svcCtx.Redis.Del(l.ctx, strconv.FormatInt(videoId, 10))

		// update count
		_, err = l.svcCtx.Redis.HIncrBy(l.ctx, VideoIDToCommentCount, strconv.FormatInt(videoId, 10), -1).Result()
		if err != nil {
			logx.Error(err)
			resp = &types.CommentResponse{StatusCode: res.InternalServerErrorCode, StatusMsg: err.Error()}
			return resp, nil
		}

		// async delete & clear cache second time
		body := fmt.Sprintf("%d-%d", commentId, videoId)
		msg := &primitive.Message{
			Topic: l.svcCtx.Config.RocketMQ.AsyncDeleteTopic,
			Body:  []byte(body),
		}
		l.svcCtx.MqProducer.SendAsync(l.ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
			if err != nil {
				logx.Error(err)
				return
			}
		}, msg)
		logx.Info("send to MQ")
		resp = &types.CommentResponse{StatusCode: res.SuccessCode, StatusMsg: "ok"}
	}
	return resp, nil
}
