package logic

import (
	"context"
	"encoding/json"
	"github.com/ev1lQuark/tiktok/service/comment/model"
	"sort"
	"time"

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
		logx.Errorf("jwt 认证失败%w", err)
		resp = &types.GetCommentListResponse{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// 校验参数
	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		logx.Errorf("参数错误%w", err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "参数错误"}
		return resp, nil
	}


	var tableComments []*model.Comment
	if n, err:= l.svcCtx.Redis.Exists(context.TODO(), strconv.FormatInt(videoId, 10)).Result(); err != nil {
		logx.Errorf("Redis error%w", err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询出错"}
		return resp, nil
	} else {
		// 没有命中缓存
		if n == 0 {
			commentQuery := l.svcCtx.Query.Comment
			tableComments, err = commentQuery.WithContext(context.TODO()).Where(commentQuery.VideoID.Eq(videoId)).Where(commentQuery.Cancel.Eq(0)).Order(commentQuery.CreatDate.Desc()).Find()
			if err != nil {
				logx.Errorf("查询错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
			tableCommentJson, err := json.Marshal(tableComments)
			if err != nil {
				logx.Errorf("json序列化错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
			_, err = l.svcCtx.Redis.Set(context.TODO(), strconv.FormatInt(videoId, 10), string(tableCommentJson), time.Duration(l.svcCtx.Config.Redis.ExpireTime) * time.Second).Result()
			if err != nil {
				logx.Errorf("Redis写入错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
		} else {
			// 命中缓存
			tableCommentJson, err := l.svcCtx.Redis.Get(context.TODO(), strconv.FormatInt(videoId, 10)).Result()
			if err != nil {
				logx.Errorf("Redis查询错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
			_, err = l.svcCtx.Redis.Expire(context.TODO(), strconv.FormatInt(videoId, 10),time.Duration(l.svcCtx.Config.Redis.ExpireTime)*time.Second).Result()
			if err != nil {
				logx.Errorf("Redis重置时间错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
			err = json.Unmarshal([]byte(tableCommentJson), &tableComments)
			if err != nil {
				logx.Errorf("Json反序列化错误:%w", err)
				resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询错误"}
				return resp, nil
			}
			sort.Slice(tableComments, func(i, j int) bool {
				return tableComments[i].CreatDate.After(tableComments[j].CreatDate)
			})
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
		userNameList, err = l.svcCtx.UserRpc.GetNames(context.TODO(), &user.IdListReq{IdList: authorIds})
		return err
	})

	// 根据userId获取本账号获赞总数
	var totalFavoriteNumList *like.GetTotalFavoriteNumReply
	eg.Go(func() error {
		var err error
		totalFavoriteNumList, err = l.svcCtx.LikeRpc.GetTotalFavoriteNum(context.TODO(), &like.GetTotalFavoriteNumReq{UserId: authorIds})
		return err
	})

	// 根据userId获取本账号喜欢（点赞）总数
	var userFavoriteCountList *like.GetFavoriteCountByUserIdReply
	eg.Go(func() error {
		var err error
		userFavoriteCountList, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserId(context.TODO(), &like.GetFavoriteCountByUserIdReq{UserId: authorIds})
		return err
	})

	// 获取work_count
	var workCount *video.VideoNumReply
	eg.Go(func() error {
		var err error
		workCount, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(context.TODO(), &video.AuthorIdReq{AuthorId: authorIds})
		return err
	})

	//错误判断
	if err := eg.Wait(); err != nil {
		logx.Errorf("Rpc调用失败%w", err)
		resp = &types.GetCommentListResponse{StatusCode: res.BadRequestCode, StatusMsg: "查询评论列表失败"}
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
				TotalFavorited:  strconv.Itoa(int(totalFavoriteNumList.Count[index])),
				WorkCount:       int(workCount.VideoNum[index]),
				FavoriteCount:   int(userFavoriteCountList.Count[index]),
			},
			Content:    value.CommentText,
			CreateDate: value.CreatDate.Format("01-02"),
		}
		commentList = append(commentList, comment)
	}
	resp = &types.GetCommentListResponse{StatusCode: res.SuccessCode, StatusMsg: "获取评论列表成功", CommentList: commentList}
	return resp, nil
}

