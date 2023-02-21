package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/types"

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

func (l *GetLikeListLogic) GetLikeList(req *types.LikeListRequest) (resp *types.LikeListResponse, err error) {
	// todo: add your logic here and delete this line
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	likeQuery := l.svcCtx.Query.Like
	//.Select(commentQuery.UserID, commentQuery.CommentText)
	//查找数据库，获取了comment表的内容,需要对result进行处理
	result, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Find()
	count, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(userId)).Count()
	// todo: 获取author信息
	//user= &types.User{
	//	ID:userId,
	//	Name :
	//	FollowCount:
	//	FollowCount
	//	IsFollow
	//}
	//commentList是接口响应返回信息
	videoList := make([]types.VideoList, count)
	for i := range result {
		//通过videoId获取当前视频受喜欢次数
		favoriteCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Count()
		//通过videoId判断用户是否对其点赞
		isF := true
		isCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Where(likeQuery.UserID.Eq(result[i].UserID)).Count()
		if isCount == 0 {
			isF = false
		}
		//对每个video进行整理,
		video := types.VideoList{
			//user还没写
			//Author:,
			//video-rpc
			//PlayURL:,
			//CoverURL :  ,
			FavoriteCount: favoriteCount,
			//comment-rpc
			//CommentCount :
			IsFavorite: isF,
			//Title      :   ,

		}
		videoList = append(videoList, video)
	}
	resp.VideoList = videoList

	return resp, nil
}
