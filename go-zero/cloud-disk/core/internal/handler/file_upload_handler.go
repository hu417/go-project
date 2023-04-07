package handler

import (
	"cloud-disk/core/helper"
	"path"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/models"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 1、根据文件hash判断文件是否存在
		// 获取请求文件信息
		file, fileheader, err := r.FormFile("file")
		if err != nil {
			return
		}
		// 根据文件size的hash判断
		b := make([]byte, fileheader.Size)
		_, err = file.Read(b)
		if err != nil {
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b))
		rp := new(models.Repository_pool)
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			return
		}
		if has {
			httpx.OkJson(w, &types.FileUploadResp{
				Identity: rp.Identity,
				Name:     rp.Name,
				Ext:      rp.Ext,
			})
			return
		}

		// 2、req参数赋值
		cospath, err := helper.CosUpload(r)
		if err != nil {
			return
		}
		req.Name = fileheader.Filename
		req.Ext = path.Ext(fileheader.Filename)
		req.Size = int(fileheader.Size)
		req.Hash = hash
		req.Path = cospath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
