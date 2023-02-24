package logic

import (
	"context"
	"fmt"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/like/setting"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdLogic {
	return &GetFavoriteCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据userId获取本账号喜欢（点赞）总数
func (l *GetFavoriteCountByUserIdLogic) GetFavoriteCountByUserId(in *like.GetFavoriteCountByUserIdReq) (*like.GetFavoriteCountByUserIdReply, error) {
	numList := make([]int64, 0, len(in.UserId))
	for _, userId := range in.UserId {
		num, err := getFavoriteCountByUserId(l.ctx, l.svcCtx, userId)
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &like.GetFavoriteCountByUserIdReply{Count: numList}, nil
}

func getFavoriteCountByUserId(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) (int64, error) {
	// 检查缓存穿透
	isPenetration, err := svcCtx.Redis.SIsMember(ctx, setting.UserIdPenetrationKey, userId).Result()
	if err != nil {
		logx.Error(err)
		return 0, err
	}
	if isPenetration {
		return 0, nil
	}

	// peek redis
	key := fmt.Sprintf(setting.UserIdKeyPattern, userId)
	res, err := svcCtx.Redis.SCard(ctx, key).Result()
	if err != nil {
		logx.Error(err)
		return 0, err
	}

	if res == 0 {
		// read from mysql
		likeQuery := svcCtx.Query.Like
		likes, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.Cancel.Eq(0)).Find()
		if err != nil {
			logx.Error(err)
			return 0, err
		}
		// cache penetration
		if len(likes) == 0 {
			svcCtx.Redis.SAdd(ctx, setting.UserIdPenetrationKey, userId)
			return 0, nil
		}

		// write to redis
		pipe := svcCtx.Redis.Pipeline()
		pipe.SRem(ctx, setting.UserIdPenetrationKey, userId)
		for _, like := range likes {
			pipe.SAdd(ctx, key, fmt.Sprintf(setting.UserIdValuePattern, like.VideoID, like.AuthorID))
		}
		pipe.Expire(ctx, key, setting.UserIdExpire)
		_, err = pipe.Exec(ctx)
		if err != nil {
			logx.Error(err)
			return 0, err
		}

		return int64(len(likes)), nil
	}
	// read from redis
	svcCtx.Redis.Expire(ctx, key, setting.UserIdExpire)
	return res, nil
}
