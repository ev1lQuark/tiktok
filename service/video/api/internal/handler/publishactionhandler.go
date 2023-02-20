package handler

import (
	"net/http"

	"github.com/ev1lQuark/tiktok/service/video/api/internal/logic"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/svc"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func publishActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		file, fileHeader, err := r.FormFile("data")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		l := logic.NewPublishActionLogic(r.Context(), svcCtx)
		resp, err := l.PublishAction(&req, &file, fileHeader)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
