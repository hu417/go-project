
$ go tool pprof ./cpu-pprof.out
(pprof) top 5  // top5 cum 、top5 flat
(pprof) traces funcxxx // 堆栈跟踪
(pprof) list funcxxx // 函数调用栈
