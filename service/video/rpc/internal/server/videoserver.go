// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package server

import (
	"context"

	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/logic"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/rpc/types/video"
)

type VideoServer struct {
	svcCtx *svc.ServiceContext
	video.UnimplementedVideoServer
}

func NewVideoServer(svcCtx *svc.ServiceContext) *VideoServer {
	return &VideoServer{
		svcCtx: svcCtx,
	}
}

func (s *VideoServer) GetVideoByVideoId(ctx context.Context, in *video.VideoIdReq) (*video.VideoInfoReply, error) {
	l := logic.NewGetVideoByVideoIdLogic(ctx, s.svcCtx)
	return l.GetVideoByVideoId(in)
}


func (s *VideoServer) GetVideoNumByAuthorId(ctx context.Context, in *video.AuthorIdReq) (*video.VideoNumReply, error) {
	l := logic.NewGetVideoNumByAuthorIdLogic(ctx, s.svcCtx)
	return l.GetVideoNumByAuthorId(in)
}

