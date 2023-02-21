package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

/*
 	Like 根据userId，videoId,actionType对视频进行点赞或者取消赞操作;
	step1: 维护Redis LikeUserId(key:strUserId),添加或者删除value:videoId,LikeVideoId(key:strVideoId),添加或者删除value:userId;z
	这里暂时用不到
	step2：更新数据库likes表;
	当前操作行为，1点赞，2取消点赞。
*/

func (l *LikeLogic) Like(req *types.LikeRequest) (resp *types.LikeResponse, err error) {
	// jwt
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		logx.Errorf("jwt 认证失败%w", err)
		resp = &types.LikeResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// 参数校验
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil || (req.ActionType != "1" && req.ActionType != "2") {
		logx.Errorf("参数错误%w")
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	likeQuery := l.svcCtx.Query.Like

	//获取Video具体信息
	VideoInfoReply, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{VideoId: []int64{videoId}})

	if err != nil {
		logx.Errorf("Rpc调用失败%w", err)
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}
	if VideoInfoReply.AuthorId[0] == 0 {
		logx.Errorf("视频不存在")
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "视频不存在"}
		return resp, nil
	}

	//点赞
	like, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.VideoID.Eq(videoId)).FirstOrCreate()
	if err != nil {
		logx.Errorf("查询数据库失败%w", err)
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "点赞失败"}
		return resp, nil
	}
	var cancel int32
	if req.ActionType == "1" {
		cancel = 0 // 点赞
	} else if req.ActionType == "2" {
		cancel = 1 // 取消点赞
	}

	_, err = likeQuery.WithContext(context.TODO()).Where(likeQuery.ID.Eq(like.ID)).UpdateSimple(likeQuery.Cancel.Value(cancel), likeQuery.AuthorID.Value(VideoInfoReply.AuthorId[0]))
	if err != nil {
		logx.Error(err)
		return &types.LikeResponse{StatusCode: res.InternalServerErrorCode, StatusMsg: "操作失败"}, nil
	}

	return &types.LikeResponse{StatusCode: res.SuccessCode, StatusMsg: "操作成功"}, nil
}
