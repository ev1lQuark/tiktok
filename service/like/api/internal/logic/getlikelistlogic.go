package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"golang.org/x/sync/errgroup"

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

// 根据userId获取点赞video列表
func (l *GetLikeListLogic) GetLikeList(req *types.LikeListRequest) (resp *types.LikeListResponse, err error) {
	// Parse jwt token
	_, err = jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.AuthFailedCode),
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误%w", err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.AuthFailedCode),
			StatusMsg:  "参数错误",
		}
		return resp, nil
	}

	likeQuery := l.svcCtx.Query.Like

	//查找数据库，获取了like表的内容,需要对result进行处理
	result, err := getLikeListByUserId(context.TODO(), l.svcCtx, userId)
	if err != nil {
		logx.Errorf("查询数据库错误%w", err)
		resp = &types.LikeListResponse{
			StatusCode: strconv.Itoa(res.InternalServerErrorCode),
			StatusMsg:  "查询数据库错误",
		}
		return resp, nil
	}

	videoId := make([]int64, 0, len(result))
	authorIds := make([]int64, 0, len(result))
	for j := range result {
		videoId = append(videoId, result[j].VideoID)
		authorIds = append(authorIds, result[j].AuthorID)
	}

	var eg errgroup.Group

	// 根据videoId获取每个视频的评论总数
	var commentReply *comment.GetComentCountByVideoIdReply
	eg.Go(func() error {
		var err error
		commentReply, err = l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &comment.GetComentCountByVideoIdReq{
			VideoId: videoId,
		})
		return err
	})

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
	var workCount *video.VideoNumReply
	eg.Go(func() error {
		var err error
		workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(l.ctx, &video.AuthorIdReq{AuthorId: authorIds})
		return err
	})

	//根据authorId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	// 错误判断
	if err := eg.Wait(); err != nil {
		logx.Errorf("调用Rpc失败%w", err)
		resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.InternalServerErrorCode), StatusMsg: "查询喜欢列表失败"}
		return resp, nil
	}

	// 通过authorId获取作者的视频喜欢数
	authorFavoriteCountList := make([]int64, 0)
	for i := range result {
		authorFavoriteCount, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(result[i].UserID)).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			logx.Errorf("数据库查询失败%w", err)
			resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.InternalServerErrorCode), StatusMsg: "查询喜欢列表失败"}
			return resp, nil
		}
		authorFavoriteCountList = append(authorFavoriteCountList, authorFavoriteCount)
	}
	// 通过authorId获取作者视频被喜欢数
	authorIsFavoriteCountList := make([]int64, 0)
	for i := range result {
		authorIsFavoriteCount, _ := likeQuery.WithContext(context.TODO()).Where(likeQuery.AuthorID.Eq(result[i].UserID)).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			logx.Errorf("数据库查询失败%w", err)
			resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.InternalServerErrorCode), StatusMsg: "查询喜欢列表失败"}
			return resp, nil
		}
		authorIsFavoriteCountList = append(authorFavoriteCountList, authorIsFavoriteCount)
	}

	videoList := make([]types.VideoList, 0, len(result))
	for i := range result {
		//通过videoId获取当前视频受喜欢次数
		favoriteCount, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			logx.Errorf("数据库查询失败%w", err)
			resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.BadRequestCode), StatusMsg: "查询喜欢列表失败"}
			return resp, nil
		}
		//通过videoId判断用户是否对其点赞
		isF := true
		isCount, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(result[i].VideoID)).Where(likeQuery.UserID.Eq(result[i].UserID)).Where(likeQuery.Cancel.Eq(0)).Count()
		if err != nil {
			logx.Errorf("数据库查询失败%w", err)
			resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.BadRequestCode), StatusMsg: "查询喜欢列表失败"}
			return resp, nil
		}
		if isCount == 0 {
			isF = false
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
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.FormatInt(authorIsFavoriteCountList[i], 10),
				WorkCount:       workCount.VideoNum[i],
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
	resp = &types.LikeListResponse{StatusCode: strconv.Itoa(res.SuccessCode), StatusMsg: "成功", VideoList: videoList}
	return resp, nil
}
