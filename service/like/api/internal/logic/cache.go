package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/setting"
	"github.com/zeromicro/go-zero/core/logx"
)

func setLike(ctx context.Context, svcCtx *svc.ServiceContext, userId, videoId, authorId int64, cancel int32) error {
	pipe := svcCtx.Redis.Pipeline()
	if cancel == 0 {
		pipe.HSet(ctx, setting.LikeMapDataKey, setting.GetLikeMapDataKey(userId, videoId), 0)
		pipe.HIncrBy(ctx, setting.LikeMapUserIdCountKey, strconv.FormatInt(userId, 10), 1)
		pipe.HIncrBy(ctx, setting.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10), 1)
		pipe.HIncrBy(ctx, setting.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10), 1)
		pipe.SAdd(ctx, setting.GetLikeSetUserIdKey(userId), videoId)
	} else {
		pipe.HSet(ctx, setting.LikeMapDataKey, setting.GetLikeMapDataKey(userId, videoId), 1)
		pipe.HIncrBy(ctx, setting.LikeMapUserIdCountKey, strconv.FormatInt(userId, 10), -1)
		pipe.HIncrBy(ctx, setting.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10), -1)
		pipe.HIncrBy(ctx, setting.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10), -1)
		pipe.SRem(ctx, setting.GetLikeSetUserIdKey(userId), videoId)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logx.Error(err)
		return err
	}

	return nil
}
