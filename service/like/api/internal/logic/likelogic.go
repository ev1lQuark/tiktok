package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/res"

	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/types"

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

	// Parse jwt token
	//userId, err := jwt.ParseUserIdFromJwtToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	//if err != nil {
	//	resp = &types.PublishActionReply{
	//		StatusCode: res.AuthFailedCode,
	//		StatusMsg:  "jwt 认证失败",
	//	}
	//	return resp, nil
	//}
	//logx.Info("userId: %v", userId)

	// 参数校验
	if req.ActionType != "1" && req.ActionType != "2" {
		resp = &types.LikeResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "点赞参数非法",
		}
		return resp, nil
	}
	if len(req.VideoId) == 0 {
		resp = &types.LikeResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	// 根据videoId和userId查询点赞记录

	return
}
