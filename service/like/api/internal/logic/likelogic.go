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
		resp = &types.LikeResponse{StatusCode: res.RemoteServiceErrorCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	//获取Video具体信息
	VideoInfoReply, err := l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{VideoId: []int64{videoId}})

	if err != nil {
		logx.Errorf("Rpc调用失败%w", err)
		resp = &types.LikeResponse{StatusCode: res.RemoteServiceErrorCode, StatusMsg: "RPC调用失败"}
		return resp, nil
	}
	if VideoInfoReply.AuthorId[0] == 0 {
		logx.Errorf("视频不存在")
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "视频不存在"}
		return resp, nil
	}

	//点赞
	var cancel int32
	if req.ActionType == "1" {
		cancel = 0 // 点赞
	} else if req.ActionType == "2" {
		cancel = 1 // 取消点赞
	}

	err = setLike(context.TODO(), l.svcCtx, userId, videoId, VideoInfoReply.AuthorId[0], cancel)
	if err != nil {
		logx.Errorf("写入数据库失败%w", err)
		resp = &types.LikeResponse{StatusCode: res.InternalServerErrorCode, StatusMsg: err.Error()}
		return resp, nil
	}

	return &types.LikeResponse{StatusCode: res.SuccessCode, StatusMsg: "操作成功"}, nil
}
