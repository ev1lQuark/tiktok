package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ev1lQuark/tiktok/service/comment/model"
	"github.com/ev1lQuark/tiktok/service/comment/pattern"
	"github.com/redis/go-redis/v9"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"

	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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
		logx.Errorf("jwt 认证失败：%w", err)
		resp = &types.GetCommentListResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// 校验参数
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误：%w", err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
		return resp, nil
	}

	var tableComments []*model.Comment
	videoIDToCommentListJSON := fmt.Sprintf(pattern.VideoIDToCommentListJSON, videoId)
	tableCommentJson, err := l.svcCtx.Redis.Get(l.ctx, videoIDToCommentListJSON).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logx.Error(err)
			resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
			return resp, nil
		} else {
			// 未命中缓存，从数据库中查询
			commentQuery := l.svcCtx.Query.Comment
			tableComments, err = commentQuery.WithContext(l.ctx).Where(commentQuery.VideoID.Eq(videoId)).Where(commentQuery.Cancel.Eq(0)).Order(commentQuery.CreatDate.Desc()).Find()
			if err != nil {
				logx.Error(err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
				return resp, nil
			}
			tableCommentJson, err := json.Marshal(tableComments)
			if err != nil {
				logx.Error(err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
				return resp, nil
			}

			_, err = l.svcCtx.Redis.Set(l.ctx, videoIDToCommentListJSON, string(tableCommentJson), time.Duration(l.svcCtx.Config.Redis.ExpireTime)*time.Second).Result()

			if err != nil {
				logx.Error(err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
				return resp, nil
			}
		}
	} else {
		_, err = l.svcCtx.Redis.Expire(l.ctx, videoIDToCommentListJSON, time.Duration(l.svcCtx.Config.Redis.ExpireTime)*time.Second).Result()
		if err != nil {
			logx.Error(err)
			resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
			return resp, nil
		}
		err = json.Unmarshal([]byte(tableCommentJson), &tableComments)
		if err != nil {
			logx.Error(err)
			resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: err.Error()}
			return resp, nil
		}
	}

	authorIds := make([]int64, 0, len(tableComments))
	for _, value := range tableComments {
		authorIds = append(authorIds, value.UserID)
	}

	var eg errgroup.Group

	//根据userId获取userName
	var userNameList *user.NameListReply
	eg.Go(func() error {
		var err error
		userNameList, err = l.svcCtx.UserRpc.GetNames(l.ctx, &user.IdListReq{IdList: authorIds})
		return err
	})

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetFavoriteCountByAuthorIdsReply
	eg.Go(func() error {
		var err error
		totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetFavoriteCountByAuthorIds(l.ctx, &like.GetFavoriteCountByAuthorIdsReq{AuthorIds: authorIds})
		return err
	})

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdsReply
	eg.Go(func() error {
		var err error
		userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserIds(l.ctx, &like.GetFavoriteCountByUserIdsReq{UserIds: authorIds})
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
		resp = &types.GetCommentListResponse{StatusCode: res.RemoteServiceErrorCode, StatusMsg: "RPC error"}
		return resp, nil
	}

	commentList := make([]types.CommentList, 0, len(tableComments))
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
				Signature:       "爱抖音，爱生活",
				TotalFavorited:  strconv.Itoa(int(totalFavoriteNumList.CountSlice[index])),
				WorkCount:       int(workCount.VideoNum[index]),
				FavoriteCount:   int(userFavoriteCountList.CountSlice[index]),
			},
			Content:    value.CommentText,
			CreateDate: value.CreatDate.Format("01-02"),
		}
		commentList = append(commentList, comment)
	}
	resp = &types.GetCommentListResponse{StatusCode: res.SuccessCode, StatusMsg: "获取评论列表成功", CommentList: commentList}
	return resp, nil
}
