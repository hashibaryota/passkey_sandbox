# ビルドステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# ビルドステージから実行ファイルとテンプレートをコピー
COPY --from=builder /main .
COPY --from=builder /app/templates ./templates/

EXPOSE 8080

CMD ["/app/main"]
