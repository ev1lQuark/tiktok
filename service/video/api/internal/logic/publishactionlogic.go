package logic

import (
	"context"
	"fmt"
	"github.com/chilts/sid"
	"github.com/minio/minio-go/v7"
	"mime/multipart"

	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/video/model"

	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq, file *multipart.File, fileHeader *multipart.FileHeader) (resp *types.PublishActionReply, err error) {
	// Parse jwt token
	//userId, err := jwt.ParseUserIdFromJwtToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	//if err != nil {
	//	resp = &types.PublishActionReply{
	//		StatusCode: res.AuthFailedCode,
	//		StatusMsg:  "jwt 认证失败",
	//	}
	//	return resp, nil
	//}
	//logx.Info("userId: %v", userId)
	userId := int64(1)
	// TODO 业务逻辑
	// uuid 视频
	videoId := sid.Id()
	videoName := videoId + fileHeader.Filename
	_, err = l.svcCtx.MinioClient.PutObject(context.TODO(), l.svcCtx.Config.Minio.VideoBucket, videoName, *file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		logx.Error("upload object error " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.BadRequestCode, StatusMsg: "上传失败"}
		return resp, nil
	}
	videoUrl := "/" + "videos/" + videoName

	// uuid 封面
	imageId := sid.Id()
	imageName := imageId + fileHeader.Filename

	//todo add video to image
	imageFile := file

	//

	_, err = l.svcCtx.MinioClient.PutObject(context.TODO(), l.svcCtx.Config.Minio.ImageBucket, imageName, *imageFile, fileHeader.Size, minio.PutObjectOptions{})

	if err != nil {
		logx.Error("upload object error " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.BadRequestCode, StatusMsg: "上传失败"}
		return resp, nil
	}
	imageUrl := "/" + "images/" + imageName

	videoQuery := l.svcCtx.Query.Video
	video := &model.Video{
		AuthorID:    userId,
		PlayURL:     videoUrl,
		CoverURL:    imageUrl,
		PublishTime: time.Now(),
		Title:       req.Title,
	}
	err = videoQuery.WithContext(context.TODO()).Create(video)
	if err != nil {
		msg := fmt.Sprintf("插入视频失败：%v", err)
		logx.Error(msg)
		resp = &types.PublishActionReply{StatusCode: res.BadRequestCode, StatusMsg: msg}
		return resp, nil
	}
	resp = &types.PublishActionReply{StatusCode: res.SuccessCode, StatusMsg: "发布成功"}
	return resp, nil
}
