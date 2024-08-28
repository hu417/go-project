package service

import (
    "fmt"
)

var Heap = new(heap)

var sum = []int{}

type heap struct {
}

func (h *heap) GetNum() {
    i := h.sum()
    for {
        i++
        sum = append(sum, i)
        if i%1000000 == 0 {
            //time.Sleep(time.Second * 1)
            fmt.Println(i)
        }
    }
}

func (h *heap) sum() int {
    num := 0
    for i := 0; i < 100000; i++ {
        num += i
    }
    return num
}