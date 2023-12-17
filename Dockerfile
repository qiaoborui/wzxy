FROM golang:latest AS builder

WORKDIR /app

# 复制应用程序代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# 第二阶段：运行阶段
FROM alpine:latest

WORKDIR /app

# 从第一阶段复制构建好的二进制文件
COPY --from=builder /app/myapp .

# 暴露应用程序需要的端口
EXPOSE 8080

# 运行应用程序
CMD ["./myapp"]