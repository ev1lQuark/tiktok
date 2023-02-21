package logic

import (
	"context"
	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/comment/model"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type CommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentLogic) Comment(req *types.GetCommentRequest) (resp *types.GetCommentResponse, err error) {
	// todo: add your logic here and delete this line
	//为1，发布评论
	actionType, err := strconv.ParseInt(req.ActionType, 10, 64)
	commentId, err := strconv.ParseInt(req.CommentId, 10, 64)
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if actionType == 1 {
		commentQuery := l.svcCtx.Query.Comment
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				break
			default:
				return nil, err
			}
		}

		videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
		logx.Info("userId: %v", userId)
		comment := &model.Comment{
			UserID:      userId,
			VideoID:     videoId,
			CommentText: req.CommentText,
			CreatDate:   time.Now(),
			Cancel:      0,
		}
		err = commentQuery.WithContext(context.TODO()).Create(comment)
		if err != nil {
			return nil, err
		}
	} else if actionType == 2 {
		commentQuery := l.svcCtx.Query.Comment
		userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				break
			default:
				return nil, err
			}
		}
		logx.Info("userId: %v", userId)

		_, err = commentQuery.WithContext(context.TODO()).Where(commentQuery.ID.Eq(commentId)).Update(commentQuery.Cancel, 1)
		if err != nil {
			return nil, err
		}
	}

	return
}
