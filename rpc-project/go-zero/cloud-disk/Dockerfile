# 多阶段构建
#构建一个 builder 镜像，目的是在其中编译出可执行文件mck
#构建时需要将此文件放到代码根目录下
FROM golang:alpine  as builder

# 环境
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.io"
#安装编译需要的环境gcc等
RUN apk add build-base

WORKDIR /build
#将上层整个文件夹拷贝到/build
ADD . /build/src
WORKDIR /build/src
#交叉编译，需要制定CGO_ENABLED=1，默认是关闭的
#去掉了调试信息 -ldflags="-s -w" 以减小镜像尺寸
RUN  GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o hello ./cmd/hello/main.go

#编译
FROM alpine
RUN apk update --no-cache && apk add --no-cache tzdata
#设置本地时区，这样我们在日志里看到的是北京时间了
ENV TZ Asia/Shanghai
WORKDIR /app
#从第一个镜像里 copy 出来可执行文件
COPY --from=builder  /build/src/hello /app/hello
#COPY --from=builder  /build/src/config/config.yml /app/config/config.yml
#COPY --from=builder  /build/src/config/alpine/libcoop_alpine.so /app/config/alpine/libcoop_alpine.so
COPY ./config/alpine/libgcc_s.so.1 /usr/lib/
#COPY ./config/alpine/libstdc++.so.6 /usr/lib/
COPY ./config/alpine/libstdc++.so.6.0.28 /usr/lib/libstdc++.so.6.0.28
RUN ln -s /usr/lib/libstdc++.so.6.0.28 /usr/lib/libstdc++.so.6
#VOLUME ["/home/tianxing/project/mck/mck-service-core/config/config.yml","/app/config/config.yml"]

# 声明服务端口
EXPOSE 8080
EXPOSE 9080

# 启动容器时运行的命令
CMD ["./hello"]