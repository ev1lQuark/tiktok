package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/rpc/types/like"
	"github.com/ev1lQuark/tiktok/service/like/setting"
	"github.com/go-redis/redis/v8"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByVideoIdLogic {
	return &GetFavoriteCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据videoId获取视频点赞总数
func (l *GetFavoriteCountByVideoIdLogic) GetFavoriteCountByVideoId(in *like.GetFavoriteCountByVideoIdReq) (*like.GetFavoriteCountByVideoIdReply, error) {
	numList := make([]int64, 0, len(in.VideoId))
	for _, videoId := range in.VideoId {
		num, err := getFavoriteCountByVideoId(l.ctx, l.svcCtx, videoId)
		if err != nil {
			return nil, err
		}
		numList = append(numList, num)
	}
	return &like.GetFavoriteCountByVideoIdReply{Count: numList}, nil
}

func getFavoriteCountByVideoId(ctx context.Context, svcCtx *svc.ServiceContext, videoId int64) (int64, error) {
	key := fmt.Sprintf(setting.VideoIdKeyPattern, videoId)
	// read from redis
	cmd := svcCtx.Redis.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存失效或过期
			// read from mysql
			likeQuery := svcCtx.Query.Like
			num, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.VideoID.Eq(videoId)).Where(likeQuery.Cancel.Eq(0)).Count()
			if err != nil {
				logx.Error(err)
				return 0, err
			}
			// write to redis
			if err := svcCtx.Redis.Set(ctx, key, num, setting.VideoIdExpire).Err(); err != nil {
				logx.Error(err)
				return 0, err
			}
			return num, nil
		}
		// unhandled redis error
		return 0, err
	}

	// read from redis
	num, err := cmd.Int64()
	if err != nil {
		logx.Error(err)
		return 0, err
	}
	return num, nil
}
