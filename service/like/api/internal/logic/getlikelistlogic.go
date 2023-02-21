package logic

import (
	"context"
	"fmt"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"golang.org/x/sync/errgroup"
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
	// Parse jwt token
	_, err = jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		resp = &types.LikeListResponse{
			StatusCode: string(res.AuthFailedCode),
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	//logx.Info("userId: %v", userId)
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误%w", err)
		resp = &types.LikeListResponse{
			StatusCode: string(res.AuthFailedCode),
			StatusMsg:  "参数错误",
		}
		return resp, nil
	}

	likeQuery := l.svcCtx.Query.Like

	//.Select(commentQuery.UserID, commentQuery.CommentText)
	//查找数据库，获取了comment表的内容,需要对result进行处理
	result, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Find()
	if err != nil {
		logx.Errorf("查询数据库错误%w", err)
		resp = &types.LikeListResponse{
			StatusCode: string(res.AuthFailedCode),
			StatusMsg:  "查询数据库错误",
		}
		return resp, nil
	}

	//user= &types.User{
	//	ID:userId,
	//	Name :
	//	FollowCount:
	//	FollowCount
	//	IsFollow
	//}

	//videoList是接口响应返回信息
	videoList := make([]types.VideoList, len(result))
	videoId := make([]int64, 0, len(result))
	authorIds := make([]int64, 0, len(result))
	for j := range result {
		videoId = append(videoId, result[j].VideoID)
		authorIds = append(authorIds, result[j].AuthorID)
	}

	var eg errgroup.Group

	// todo:comment-rpc:
	// 根据videoId获取每个视频的评论总数
	var commentReply *comment.GetComentCountByVideoIdReply
	eg.Go(func() error {
		var err error
		commentReply, err = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{
			VideoId: videoId,
		})
		return err
	})

	// todo:video-rpc
	/*
		AuthorId:    authorIdList,
		PlayUrl:     playUrlList,
		CoverUrl:    coverUrlList,
		PublishTime: publishTimeList,
		Title:       tileList,
	*/
	var VideoInfoReply *video.VideoInfoReply
	//获取Video具体信息
	eg.Go(func() error {
		var err error
		VideoInfoReply, err = l.svcCtx.VideoRpc.GetVideoByVideoId(l.ctx, &video.VideoIdReq{
			VideoId: videoId,
		})
		return err
	})

	// 获取workCount
	workCount := make([]int64, 0, len(result))
	var videoNumReply *video.VideoNumReply
	eg.Go(func() error {
		var err error
		videoNumReply, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{
			AuthorId: authorIds,
		})
		return err
	})
	for _, value := range videoNumReply.VideoNum {
		workCount = append(workCount, value)
	}

	//CoverURL，
	//title

	// todo:user-rpc
	//根据authorId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	//错误判断
	if err := eg.Wait(); err != nil {
		msg := fmt.Sprintf("調用Rpc失敗%v", err)
		logx.Error(msg)
		resp = &types.LikeListResponse{StatusCode: string(res.BadRequestCode), StatusMsg: msg}
		return resp, nil
	}

	//通过authorId获取作者的视频喜欢数
	authorFavoriteCountList := make([]int64, 0, len(result))
	for i := range result {
		authorFavoriteCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(result[i].UserID)).Count()
		authorFavoriteCountList = append(authorFavoriteCountList, authorFavoriteCount)
	}
	//通过authorId获取作者视频被喜欢数
	authorisFavoritedCountList := make([]int64, 0, len(result))
	for i := range result {
		authorisFFavoritedCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.AuthorID.Eq(result[i].UserID)).Count()
		authorisFavoritedCountList = append(authorFavoriteCountList, authorisFFavoritedCount)
	}
	for i := range result {
		//通过videoId获取当前视频受喜欢次数
		favoriteCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Count()
		//通过videoId判断用户是否对其点赞
		isF := true
		isCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Where(likeQuery.UserID.Eq(result[i].UserID)).Count()
		if isCount == 0 {
			isF = false
		}
		if err != nil {
			return nil, err
		}

		//对每个video进行整理,
		videoSingle := types.VideoList{
			ID: videoId[i],
			Author: types.Author{
				ID:              authorIds[i],
				Name:            userNameList.NameList[i],
				FollowCount:     0,
				FollowerCount:   0,
				IsFollow:        false,
				Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
				Signature:       "愛抖音，爱生活",
				TotalFavorited:  strconv.FormatInt(authorisFavoritedCountList[i], 10),
				WorkCount:       workCount[i],
				FavoriteCount:   authorFavoriteCountList[i],
			},
			PlayURL:       VideoInfoReply.PlayUrl[i],
			CoverURL:      VideoInfoReply.CoverUrl[i],
			FavoriteCount: favoriteCount,
			CommentCount:  commentReply.Count[i],
			IsFavorite:    isF,
			Title:         VideoInfoReply.Title[i],
		}
		videoList = append(videoList, videoSingle)
	}

	resp = &types.LikeListResponse{StatusCode: string(res.SuccessCode), StatusMsg: "成功", VideoList: videoList}
	return resp, nil
}
