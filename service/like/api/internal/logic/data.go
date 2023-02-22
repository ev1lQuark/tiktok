package logic

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/model"
	"github.com/zeromicro/go-zero/core/logx"
)

func readLike(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) ([]*model.Like, error) {
	key := fmt.Sprintf("like:%d", userId)
	res, err := svcCtx.Redis.SMembers(context.TODO(), key).Result()
	if err != nil {
		// redis error
		logx.Error(err)
		return nil, err
	}
	if len(res) == 0 {
		// 缓存失效或过期
		logx.Debug("read from mysql")
		// read from mysql
		likeQuery := svcCtx.Query.Like
		result, err := likeQuery.WithContext(context.TODO()).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.Cancel.Eq(0)).Find()
		if err != nil {
			logx.Errorf("查询数据库错误%w", err)
			return nil, err
		}
		// write to redis
		for _, v := range result {
			value := fmt.Sprintf("%d-%d", v.VideoID, v.AuthorID)
			if err := svcCtx.Redis.SAdd(context.TODO(), key, value).Err(); err != nil {
				logx.Errorf("写入redis错误%w", err)
				return nil, err
			}
		}
		return result, nil
	}

	// read from redis
	logx.Debug("read from redis")
	likeList := make([]*model.Like, 0, len(res))
	for _, v := range res {
		like := &model.Like{UserID: userId}
		_, err := fmt.Sscanf(v, "%d-%d", &like.VideoID, &like.AuthorID)
		if err != nil {
			logx.Errorf("redis数据格式错误%w", err)
			return nil, err
		}
		likeList = append(likeList, like)
	}
	return likeList, nil
}

func writeLike(ctx context.Context, svcCtx *svc.ServiceContext, userId, videoId, authorId int64, cancel int32) error {
	key := fmt.Sprintf("like:%d", userId)
	value := fmt.Sprintf("%d-%d", videoId, authorId)

	// 修改 redis
	if cancel == 0 {
		// 点赞
		err := svcCtx.Redis.SAdd(ctx, key, value).Err()
		if err != nil {
			logx.Error(err)
			return err
		}
	} else {
		// 取消点赞
		err := svcCtx.Redis.SRem(ctx, key, value).Err()
		if err != nil {
			logx.Error(err)
			return err
		}
	}

	// 修改 mysql
	body := fmt.Sprintf("like:%d-%d-%d-%d", userId, videoId, authorId, cancel)
	msg := &primitive.Message{
		Topic: svcCtx.Config.RocketMQ.Topic,
		Body:  []byte(body),
	}

	// 发送消息到 MQ
	svcCtx.MqProducer.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			logx.Error(err)
			return
		}
	}, msg)
	logx.Info("send to MQ")

	return nil
}
