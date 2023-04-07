package handler

import (
	"book/service/shorturl/api/internal/logic"
	"book/service/shorturl/api/internal/svc"
	"book/service/shorturl/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ShortenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShortenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShortenLogic(r.Context(), svcCtx)
		resp, err := l.Shorten(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
