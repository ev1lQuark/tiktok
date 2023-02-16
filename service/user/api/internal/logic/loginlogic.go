package logic

import (
	"context"
	"strings"
	"time"

	"github.com/ev1lQuark/tiktok/common/encrypt"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/types"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginOrRegisterReq) (resp *types.LoginOrRegisterReply, err error) {
	// 参数校验
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		resp = &types.LoginOrRegisterReply{Code: res.DefaultErrorCode, Message: "参数错误"}
		return resp, nil
	}

	userQuery := l.svcCtx.Query.User
	// 判断用户是否存在
	user, err := userQuery.WithContext(context.TODO()).Where(userQuery.Name.Eq(req.Username)).Take()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			resp = &types.LoginOrRegisterReply{Code: res.DefaultErrorCode, Message: "用户不存在"}
			return resp, nil
		default:
			return nil, err
		}
	}

	// 校验密码
	p := encrypt.Sha256Encrypt(req.Password)
	if user.Password != p {
		resp = &types.LoginOrRegisterReply{Code: res.DefaultErrorCode, Message: "密码错误"}
		return resp, nil
	}

	// 生成token
	token, err := jwt.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.ID)
	if err != nil {
		return nil, err
	}

	return &types.LoginOrRegisterReply{Code: res.SuccessCode, Message: "登录成功", UserId: user.ID, Token: token}, nil
}
