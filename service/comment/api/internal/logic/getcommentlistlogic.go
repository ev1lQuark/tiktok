package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
	"strconv"
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
	_, err = jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		logx.Error(err)
		resp = &types.GetCommentListResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}
	// 校验参数
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		logx.Error(err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}

	//
	commentQuery := l.svcCtx.Query.Comment

	// 查找数据库，获取了comment表的内容,需要对result进行处理
	tableComments, err := commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Where(commentQuery.Cancel.Eq(1)).Order(commentQuery.CreatDate.Desc()).Find()
	if err != nil {
		msg := "查询出错"
		logx.Error(msg + err.Error())
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}
	//
	authorIds := make([]int64, 0, len(tableComments))
	videoIds := make([]int64, 0, len(tableComments))

	for _, value := range tableComments {
		authorIds = append(authorIds, value.UserID)
		videoIds = append(videoIds, value.VideoID)
	}
	//

	var eg errgroup.Group

	//根据userId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetTotalFavoriteNumReply

	eg.Go(func() error {
		var err error
		totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetTotalFavoriteNum(l.ctx, &like.GetTotalFavoriteNumReq{UserId: authorIds})
		return err
	})

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
	eg.Go(func() error {
		var err error
		userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(l.ctx, &like.GetFavoriteCountByUserIdReq{UserId: authorIds})
		return err
	})

	// 获取work_count
	var workCount *video.VideoNumReply
	eg.Go(func() error {
		var err error
		workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{AuthorId: authorIds})
		return err
	})

	//错误判断
	if err := eg.Wait(); err != nil {
		logx.Errorf("Rpc调用失败%w", err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询评论列表失败"}
		return resp, nil
	}

	commentList := make([]types.CommentList, len(tableComments))
	for index, value := range tableComments {
		//对每个评论进行整理
		comment := types.CommentList{
			ID: value.ID,
			User: types.User{
				ID:              int(authorIds[index]),
				Name:            userNameList.NameList[index],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "愛抖音，爱生活",
				TotalFavorited:  strconv.Itoa(int(totalFavoriteNumList.Count[index])),
				WorkCount:       int(workCount.VideoNum[index]),
				FavoriteCount:   int(userFavoriteCountList.Count[index]),
			},
			Content:    value.CommentText,
			CreateDate: value.CreatDate.String(),
		}
		commentList = append(commentList, comment)
	}

	resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "获取评论列表成功", CommentList: commentList}
	return resp, nil
}
