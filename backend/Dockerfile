# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk --no-cache add git

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o remotegpu ./cmd/

# 运行阶段
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/remotegpu .
COPY --from=builder /app/config/config.yaml.example ./config/config.yaml.example

EXPOSE 8080

CMD ["./remotegpu", "server"]
