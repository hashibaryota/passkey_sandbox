# ビルドステージ
FROM public.ecr.aws/docker/library/golang:1.25.1-alpine3.22 AS builder

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# 実行ステージ
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# 実行ファイルとテンプレートをコピー
COPY --from=builder /app/server .
COPY ./templates/ ./templates/

CMD ["/app/server"]
