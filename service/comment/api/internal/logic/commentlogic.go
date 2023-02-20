package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/comment/model"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/zeromicro/go-zero/core/logx"
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
		resp = &types.GetCommentResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// 参数校验
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	actionType, err := strconv.ParseInt(req.ActionType, 10, 64)
	if err != nil || (actionType != int64(1) && actionType != int64(2)) {
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	commentId, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	// 判断videoId是否存在
	videoInfo, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{VideoId: []int64{videoId}})

	if err != nil {
		resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "rpc调用失败" + err.Error()}
		return resp, nil
	}

	if videoInfo.AuthorId[0] == 0 {
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
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "发布评论失败"}
			return resp, nil
		}
	} else if actionType == 2 {
		commentQuery := l.svcCtx.Query.Comment
		_, err = commentQuery.WithContext(context.TODO()).Where(commentQuery.ID.Eq(commentId)).Update(commentQuery.Cancel, 1)
		if err != nil {
			resp = &types.GetCommentResponse{StatusCode: res.BadRequestCode, StatusMsg: "删除评论失败"}
			return resp, nil
		}
	}

	return
}
