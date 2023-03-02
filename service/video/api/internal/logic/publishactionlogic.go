package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/chilts/sid"
	"github.com/minio/minio-go/v7"

	"github.com/ev1lQuark/tiktok/common/jwt"
	"github.com/ev1lQuark/tiktok/common/res"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/ev1lQuark/tiktok/service/video/model"
	"github.com/ev1lQuark/tiktok/service/video/pattern"

	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/zeromicro/go-zero/core/logx"
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
	userId, err := jwt.GetUserId(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	if err != nil {
		logx.Error(err)
		resp = &types.PublishActionReply{
			StatusCode: res.AuthFailedCode,
			StatusMsg:  "jwt 认证失败",
		}
		return resp, nil
	}

	// uuid 视频
	videoId := sid.Id()
	videoName := videoId + fileHeader.Filename
	// 上传视频
	_, err = l.svcCtx.MinioClient.PutObject(context.TODO(), l.svcCtx.Config.Minio.VideoBucket, videoName, *file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		logx.Error("upload object error " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "视频上传失败"}
		return resp, nil
	}
	videoUrl := "/" + "videos/" + videoName

	// uuid 封面
	imageId := sid.Id()
	imageName := imageId + fileHeader.Filename[:strings.LastIndex(fileHeader.Filename, ".")] + ".jpg"

	// 视频抽帧作为封面
	imageFile, err := readFrameFromVideo(fmt.Sprintf("http://%s%s", l.svcCtx.Config.Minio.Endpoint, videoUrl))
	if err != nil {
		logx.Error("read frame from video error " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "视频截取封面出错"}
		return resp, nil
	}

	// 上传封面图
	_, err = l.svcCtx.MinioClient.PutObject(context.TODO(), l.svcCtx.Config.Minio.ImageBucket, imageName, imageFile, -1, minio.PutObjectOptions{})
	if err != nil {
		logx.Error("upload object error " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "封面图上传失败"}
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

	videoJSON, err := json.Marshal(video)
	if err != nil {
		logx.Error("JSON序列化出错:%w " + err.Error())
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "失败"}
		return resp, nil
	}

	l1 := struct {
		Score  float64
		Member interface{}
	}{
		float64(video.PublishTime.Unix()),
		string(videoJSON),
	}

	_, err = l.svcCtx.Redis.ZAdd(context.TODO(), pattern.VideoDataListJSON, l1).Result()
	if err != nil {
		logx.Error("Redis插入新视频出错%w", err)
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "失败"}
		return resp, nil
	}

	// Redis的map进行increase若不存在key，默认为0
	_, err = l.svcCtx.Redis.HIncrBy(context.TODO(), pattern.AuthorIdToWorkCount, strconv.FormatInt(userId, 10), 1).Result()
	if err != nil {
		logx.Error("Redis增加对应workCount失败%w", err)
		resp = &types.PublishActionReply{StatusCode: res.InternalServerErrorCode, StatusMsg: "失败"}
		return resp, nil
	}

	resp = &types.PublishActionReply{StatusCode: res.SuccessCode, StatusMsg: "发布成功"}
	return resp, nil
}

func readFrameFromVideo(inputUrl string) (frame io.Reader, err error) {
	outBuf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(inputUrl).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).WithOutput(outBuf).Run()
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	frame = outBuf
	return frame, err
}
