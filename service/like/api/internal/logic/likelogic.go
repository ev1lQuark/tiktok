package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/service/Like/model"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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
	// todo: add your logic here and delete this line
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	logx.Info("userId: %v", userId)
	likeQuery := l.svcCtx.Query.Like
	//actionType, err := strconv.ParseInt(req.ActionType, 10, 32)
	//点赞
	if req.ActionType == "1" {
		like := &model.Like{
			UserID:  userId,
			VideoID: videoId,
			Cancel:  0,
			// todo:通过video rpc:由videoid获取authorid
			//AuthorID:
		}
		err := likeQuery.WithContext(context.TODO()).Create(like)
		if err != nil {
			return nil, err
		} //取消点赞  通过userid和videoid
	} else if req.ActionType == "2" {
		_, err = likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.VideoID.Eq(videoId)).Update(likeQuery.Cancel, 1)
		if err != nil {
			return nil, err
		}
	}
	return
}
