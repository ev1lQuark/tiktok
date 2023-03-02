package logic

import (
	"context"
	"strconv"

	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/pattern"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

func setLike(ctx context.Context, svcCtx *svc.ServiceContext, userId, videoId, authorId int64, cancel int32) error {
	s, err := svcCtx.Redis.HGet(ctx, pattern.LikeMapDataKey, pattern.GetLikeMapDataKey(userId, videoId)).Result()
	if err != nil {
		if err != redis.Nil {
			logx.Error(err)
			return err
		}
	} else {
		if s == "0" && cancel == 0 || s == "1" && cancel == 1 {
			return nil
		}
	}
	pipe := svcCtx.Redis.Pipeline()
	if cancel == 0 {
		pipe.HSet(ctx, pattern.LikeMapDataKey, pattern.GetLikeMapDataKey(userId, videoId), 0)
		pipe.HIncrBy(ctx, pattern.LikeMapUserIdCountKey, strconv.FormatInt(userId, 10), 1)
		pipe.HIncrBy(ctx, pattern.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10), 1)
		pipe.HIncrBy(ctx, pattern.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10), 1)
		pipe.SAdd(ctx, pattern.GetLikeSetUserIdKey(userId), videoId)
	} else {
		pipe.HSet(ctx, pattern.LikeMapDataKey, pattern.GetLikeMapDataKey(userId, videoId), 1)
		pipe.HIncrBy(ctx, pattern.LikeMapUserIdCountKey, strconv.FormatInt(userId, 10), -1)
		pipe.HIncrBy(ctx, pattern.LikeMapVideoIdCountKey, strconv.FormatInt(videoId, 10), -1)
		pipe.HIncrBy(ctx, pattern.LikeMapAuthorIdCountKey, strconv.FormatInt(authorId, 10), -1)
		pipe.SRem(ctx, pattern.GetLikeSetUserIdKey(userId), videoId)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		logx.Error(err)
		return err
	}

	return nil
}
