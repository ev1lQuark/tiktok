package logic

import (
	"context"
	"fmt"
	"github.com/ev1lQuark/tiktok/common/config"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"log"
	"strconv"

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
	//登录校验
	jwt.ParseUserIdFromJwtToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
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
	//查找last date最近视屏
	videoQuery := l.svcCtx.Query.Video

	tableVideos, err := videoQuery.WithContext(context.TODO()).Where(videoQuery.AuthorID.Eq(userId)).Order(videoQuery.PublishTime.Desc()).Find()

	if err != nil {
		log.Printf("获取用户的视频发布列表失败：%v", err)
		resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: "获取用户的视频发布列表失败"}
		return resp, nil
	}
	log.Printf("获取用户的视频发布列表成功：%v", tableVideos)

	videos := make([]types.VideoList, 0, config.VideoCount)
	for _, value := range tableVideos {
		videos = append(videos, types.VideoList{
			ID: int(value.ID),
			//todo add rpc Author
			//Author:
			PlayURL:  value.PlayURL,
			CoverURL: value.CoverURL,
			//todo add rpc favorite_count comment_count
			//FavoriteCount: val.
			//CommentCount:
			//IsFavorite:
			Title: value.Title,
		})
		fmt.Println(value.PublishTime)
	}
	resp = &types.PublishListReply{StatusCode: res.SuccessCode, StatusMsg: "成功", VideoList: videos}
	return resp, nil
}
