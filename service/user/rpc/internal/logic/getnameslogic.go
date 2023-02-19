package logic

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/user/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNamesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNamesLogic {
	return &GetNamesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNamesLogic) GetNames(in *user.IdListReq) (*user.NameListReply, error) {
	u := l.svcCtx.Query.User

	userList, err := u.WithContext(context.TODO()).Select(u.Name).Where(u.ID.In(in.IdList...)).Find()
	if err != nil {
		return nil, err
	}

	nameList := make([]string, 0, len(userList))
	for _, u := range userList {
		nameList = append(nameList, u.Name)
	}

	return &user.NameListReply{
		NameList: nameList,
	}, nil
}
