package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/types"

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

	logx.Debug(user)
	// TODO 完成UserInfo剩余逻辑

	return
}
