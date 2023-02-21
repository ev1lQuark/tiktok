package logic

import (
	"context"
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
	"strconv"
	"time"
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

func (l *CommentLogic) Comment(req *types.GetCommentRequest) (resp *types.GetCommentResponse, err error) {
	// Parse jwt token
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		logx.Errorf("jwt 认证失败%w", err)
		resp = &types.GetCommentResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// 参数校验
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误%w", err)
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	actionType, err := strconv.ParseInt(req.ActionType, 10, 64)
	if err != nil || (actionType != int64(1) && actionType != int64(2)) {
		logx.Errorf("参数错误%w", err)
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	// 判断videoId是否存在
	videoInfo, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{VideoId: []int64{videoId}})

	if err != nil {
		logx.Errorf("Rpc调用失败%w", err)
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	if videoInfo.AuthorId[0] == 0 {
		logx.Errorf("视频不存在")
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "视频不存在"}
		return resp, nil
	}

	//为1，发布评论
	if actionType == 1 {
		commentQuery := l.svcCtx.Query.Comment
		comment := &model.Comment{
			UserID:      userId,
			VideoID:     videoId,
			CommentText: req.CommentText,
			CreatDate:   time.Now(),
			Cancel:      0,
		}
		err = commentQuery.WithContext(context.TODO()).Create(comment)
		if err != nil {
			logx.Errorf("发布评论失败%w", err)
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "发布评论失败"}
			return resp, nil
		}

		var eg errgroup.Group

		//根据userId获取userName
		var userNameList *user.NameListReply
		eg.Go(func() error {
			var err error
			userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: []int64{userId}})
			return err
		})

		// 根据userId获取本账号获赞总数
		var totalFavoriteNumList *like.GetTotalFavoriteNumReply

		eg.Go(func() error {
			var err error
			totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetTotalFavoriteNum(l.ctx, &like.GetTotalFavoriteNumReq{UserId: []int64{userId}})
			return err
		})

		// 根据userId获取本账号喜欢（点赞）总数
		var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
		eg.Go(func() error {
			var err error
			userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(l.ctx, &like.GetFavoriteCountByUserIdReq{UserId: []int64{userId}})
			return err
		})

		// 获取work_count
		var workCount *video.VideoNumReply
		eg.Go(func() error {
			var err error
			workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{AuthorId: []int64{userId}})
			return err
		})

		//错误判断
		if err := eg.Wait(); err != nil {
			logx.Errorf("Rpc调用失败%w", err)
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "获取评论详细信息失败"}
			return resp, nil
		}

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
				Signature:       "愛抖音，爱生活",
				TotalFavorited:  strconv.Itoa(int(totalFavoriteNumList.Count[0])),
				WorkCount:       int(workCount.VideoNum[0]),
				FavoriteCount:   int(userFavoriteCountList.Count[0]),
			},
			Content:    comment.CommentText,
			CreateDate: comment.CreatDate.String(),
		}
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "发布评论成功", Comment: commentInfo}
	} else if actionType == 2 {
		// 参数校验
		commentId, err := strconv.ParseInt(req.CommentId, 10, 64)
		if err != nil {
			logx.Errorf("参数错误%w", err)
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
			return resp, nil
		}
		commentQuery := l.svcCtx.Query.Comment
		_, err = commentQuery.WithContext(context.TODO()).Where(commentQuery.ID.Eq(commentId)).Update(commentQuery.Cancel, 1)
		if err != nil {
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "删除评论失败"}
			return resp, nil
		}
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "删除评论成功"}
	}
	return resp, nil
}
