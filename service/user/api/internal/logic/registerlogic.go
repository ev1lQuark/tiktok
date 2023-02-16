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
	"github.com/ev1lQuark/tiktok/service/user/model"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.LoginOrRegisterReq) (resp *types.LoginOrRegisterReply, err error) {
	// 参数校验
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		resp = &types.LoginOrRegisterReply{Code: res.DefaultErrorCode, Message: "参数错误"}
		return resp, nil
	}

	userQuery := l.svcCtx.Query.User
	// 判断用户是否存在
	_, err = userQuery.WithContext(context.TODO()).Where(userQuery.Name.Eq(req.Username)).Take()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			break
		default:
			return nil, err
		}
	} else {
		resp = &types.LoginOrRegisterReply{Code: res.DefaultErrorCode, Message: "用户已存在"}
		return resp, nil
	}

	// 创建用户
	user := &model.User{
		Name:     req.Username,
		Password: encrypt.Sha256Encrypt(req.Password),
	}
	err = userQuery.WithContext(context.TODO()).Create(user)

	if err != nil {
		return nil, err
	}

	// 生成jwt
	token, err := jwt.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.ID)
	if err != nil {
		return nil, err
	}

	return &types.LoginOrRegisterReply{Code: res.SuccessCode, Message: "注册成功", UserId: user.ID, Token: token}, nil
}
