# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o remotegpu-proxy ./cmd/

# 运行阶段
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/remotegpu-proxy .

CMD ["./remotegpu-proxy"]
