package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
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
	// Parse jwt token
	//userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	//if err != nil {
	//	logx.Error(err)
	//	resp = &types.GetCommentListResponse{
	//		StatusCode: res.AuthFailedCode,
	//		StatusMsg:  "jwt 认证失败",
	//	}
	//	return resp, nil
	//}
	// 校验参数
	//videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	//if err != nil {
	//	logx.Error(err)
	//	resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
	//	return resp, nil
	//}

	//
	//commentQuery := l.svcCtx.Query.Comment
	//.Select(commentQuery.UserID, commentQuery.CommentText)
	//查找数据库，获取了comment表的内容,需要对result进行处理
	//tableComments, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Where(commentQuery.Cancel.Eq(1)).Order(commentQuery.CreatDate.Desc()).Find()
	//if err != nil {
	//	msg := "查询出错"
	//	logx.Error(msg + err.Error())
	//	resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: msg}
	//	return resp, nil
	//}
	//// todo: 获取user信息
	//authorIds := make([]int64, 0, len(tableComments))
	//videoIds := make([]int64, 0, len(tableComments))
	//
	//for _, value := range tableComments {
	//	authorIds = append(authorIds, value.UserID)
	//	videoIds = append(videoIds, value.VideoID)
	//}
	//
	//var eg errgroup.Group
	//
	////根据userId获取userName
	//var userNameList *user.NameListReply
	//eg.Go(func() error {
	//	var err error
	//	userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
	//	return err
	//})
	//
	//// 根据userId获取本账号获赞总数
	//var totalFavoriteNumList *like.GetTotalFavoriteNumReply
	//
	//eg.Go(func() error {
	//	var err error
	//	totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetTotalFavoriteNum(l.ctx, &like.GetTotalFavoriteNumReq{UserId: authorIds})
	//	return err
	//})
	//
	//// 根据userId获取本账号喜欢（点赞）总数
	//var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
	//eg.Go(func() error {
	//	var err error
	//	userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(l.ctx, &like.GetFavoriteCountByUserIdReq{UserId: authorIds})
	//	return err
	//})
	//
	////错误判断
	//if err := eg.Wait(); err != nil {
	//	msg := fmt.Sprintf("調用Rpc失敗%v", err)
	//	logx.Error(msg)
	//	resp = &types.PublishListReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
	//	return resp, nil
	//}

	//user= &types.User{
	//	ID:userId,
	//	Name :
	//	FollowCount:
	//	FollowCount
	//	IsFollow
	//}
	//commentList是接口响应返回信息
	//commentList := make([]types.CommentList, count)
	//for i := range result {
	//	//对每个评论进行整理
	//	comment := types.CommentList{
	//		ID: result[i].ID,
	//		//user还没写
	//		//User:,
	//		Content:    result[i].CommentText,
	//		CreateDate: result[i].CreatDate.String(),
	//	}
	//	commentList = append(commentList, comment)
	//}
	//resp.CommentList = commentList

	return resp, nil
}
