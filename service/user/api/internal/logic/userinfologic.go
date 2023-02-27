package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
	"golang.org/x/sync/errgroup"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoReply, err error) {
	if !jwt.Verify(l.svcCtx.Config.Auth.AccessSecret, req.Token) {
		return &types.UserInfoReply{
			StautsCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}, nil
	}

	u := l.svcCtx.Query.User
	user, err := u.WithContext(context.TODO()).Where(u.ID.Eq(req.UserId)).First()
	if err != nil {
		resp = &types.UserInfoReply{
			StautsCode: res.BadRequestCode,
			StatusMsg:  "user not found",
		}
		return resp, nil
	}

	userId := user.ID

	var eg errgroup.Group

	var videoNumReply *video.VideoNumReply
	eg.Go(func() error {
		var err error
		videoNumReply, err = l.svcCtx.VideoRpc.GetVideoNumByAuthorId(context.TODO(), &video.AuthorIdReq{AuthorId: []int64{userId}})
		return err
	})

	var totalFavoriteNumReply *like.GetFavoriteCountByAuthorIdsReply
	eg.Go(func() error {
		var err error
		totalFavoriteNumReply, err = l.svcCtx.LikeRpc.GetFavoriteCountByAuthorIds(context.TODO(), &like.GetFavoriteCountByAuthorIdsReq{AuthorIds: []int64{userId}})
		return err
	})

	var favoriteCountReply *like.GetFavoriteCountByUserIdsReply
	eg.Go(func() error {
		var err error
		favoriteCountReply, err = l.svcCtx.LikeRpc.GetFavoriteCountByUserIds(context.TODO(), &like.GetFavoriteCountByUserIdsReq{UserIds: []int64{userId}})
		return err
	})

	// 错误判断
	if err := eg.Wait(); err != nil {
		logx.Errorf("调用Rpc失败%w", err)
		resp = &types.UserInfoReply{StautsCode: res.InternalServerErrorCode, StatusMsg: "调用Rpc失败"}
		return resp, nil
	}

	return &types.UserInfoReply{
		StautsCode: res.SuccessCode,
		StatusMsg:  "success",
		User: &types.User{
			ID:              userId,
			Name:            user.Name,
			FollowCount:     0,
			FollowerCount:   0,
			IsFollow:        false,
			Avatar:          "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
			BackgroundImage: "https://inews.gtimg.com/newsapp_bt/0/13352207849/1000",
			Signature:       "爱抖音，爱生活",
			TotalFavorited:  strconv.Itoa(int(totalFavoriteNumReply.CountSlice[0])),
			WorkCount:       videoNumReply.VideoNum[0],
			FavoriteCount:   favoriteCountReply.CountSlice[0],
		},
	}, nil
}
