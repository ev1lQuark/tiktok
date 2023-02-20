package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 通过videoID查找所有评论
func (l *GetCommentListLogic) GetCommentList(req *types.GetCommentListRequest) (resp *types.GetCommentListResponse, err error) {
	// todo: add your logic here and delete this line
	//userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	commentQuery := l.svcCtx.Query.Comment
	//.Select(commentQuery.UserID, commentQuery.CommentText)
	//查找数据库，获取了comment表的内容,需要对result进行处理
	result, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Order(commentQuery.CreatDate).Find()
	count, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Count()
	// todo: 获取user信息

	//user= &types.User{
	//	ID:userId,
	//	Name :
	//	FollowCount:
	//	FollowCount
	//	IsFollow
	//}
	//commentList是接口响应返回信息
	commentList := make([]types.CommentList, count)
	for i := range result {
		//对每个评论进行整理
		comment := types.CommentList{
			ID: result[i].ID,
			//user还没写
			//User:,
			Content:    result[i].CommentText,
			CreateDate: result[i].CreatDate.String(),
		}
		commentList = append(commentList, comment)
	}
	resp.CommentList = commentList

	return resp, nil
}
