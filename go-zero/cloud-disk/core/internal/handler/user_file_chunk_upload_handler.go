package handler

import (
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserFileChunkUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileChunkUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 对请求参数的判断
		if r.PostForm.Get("key") == "" || r.PostForm.Get("part_number") == "" || r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("key or upload_id or post_id is empty"))
			return
		}

		// 获取etag
		etag, err := helper.CosPartUpload(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserFileChunkUploadLogic(r.Context(), svcCtx)
		resp, err := l.UserFileChunkUpload(&req, etag)
		// resp = new(types.UserFileChunkUploadResp)
		// resp.Etag = etag
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
