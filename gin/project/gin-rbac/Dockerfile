# 第一阶段：构建阶段
FROM golang:1.22.2-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY 代理
ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目代码
COPY . .

# 构建可执行文件
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/service/main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件和 Swagger 文档
COPY --from=builder /app/main /root
COPY --from=builder /app/docs /root/docs
COPY --from=builder /app/settings /root/settings

# 确保文件具有执行权限
RUN chmod +x /root/main

# 暴露端口
EXPOSE 8081

# 运行可执行文件
CMD ["sh", "-c", "./main"]