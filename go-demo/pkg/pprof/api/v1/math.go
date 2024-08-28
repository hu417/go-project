package v1

import (
    "github.com/gin-gonic/gin"
    "pprof/service"
)

var Heap = new(heap)

type heap struct {
}

func (h *heap) Log(ctx *gin.Context) {
    service.Heap.GetNum()
}
