package logic

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/like/model"
	"github.com/ev1lQuark/tiktok/service/like/setting"
	"github.com/zeromicro/go-zero/core/logx"
)

func getLikeListByUserId(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) ([]*model.Like, error) {
	// 防止缓存穿透
	isPenetration, err := svcCtx.Redis.SIsMember(ctx, setting.UserIdPenetrationKey, userId).Result()
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if isPenetration {
		return nil, nil
	}

	// peek redis
	key := fmt.Sprintf(setting.UserIdKeyPattern, userId)
	res, err := svcCtx.Redis.SMembers(ctx, key).Result()
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if len(res) == 0 {
		likeQuery := svcCtx.Query.Like
		result, err := likeQuery.WithContext(ctx).Where(likeQuery.UserID.Eq(userId)).Where(likeQuery.Cancel.Eq(0)).Find()
		if err != nil {
			logx.Errorf("查询数据库错误%w", err)
			return nil, err
		}
		if len(result) == 0 {
			svcCtx.Redis.SAdd(ctx, setting.UserIdPenetrationKey, userId) // 加入缓存穿透set
			return nil, nil
		} else {
			// write to redis
			for _, v := range result {
				value := fmt.Sprintf(setting.UserIdValuePattern, v.VideoID, v.AuthorID)
				if err := svcCtx.Redis.SAdd(ctx, key, value).Err(); err != nil {
					logx.Errorf("写入redis错误%w", err)
					return nil, err
				}
			}
			refreshExpire(ctx, svcCtx, key)
			return result, nil
		}

	}

	// read from redis
	refreshExpire(ctx, svcCtx, key)
	likeList := make([]*model.Like, 0, len(res))
	for _, v := range res {
		like := &model.Like{UserID: userId}
		_, err := fmt.Sscanf(v, setting.UserIdValuePattern, &like.VideoID, &like.AuthorID)
		if err != nil {
			logx.Errorf("redis数据格式错误%w", err)
			return nil, err
		}
		likeList = append(likeList, like)
	}
	return likeList, nil
}

func setLike(ctx context.Context, svcCtx *svc.ServiceContext, userId, videoId, authorId int64, cancel int32) error {
	body := fmt.Sprintf(svc.MsgPattern, userId, videoId, authorId, cancel)
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

	return nil
}

func refreshExpire(ctx context.Context, svcCtx *svc.ServiceContext, key string) {
	if ok := svcCtx.Redis.Expire(ctx, key, setting.UserIdExpire); !ok.Val() {
		logx.Error("设置过期时间失败")
	}
}
