FROM golang:1.17-alpine as builder
RUN apk update && apk add --no-cache ca-certificates
WORKDIR /app
COPY go.* .
RUN --mount=type=cache,target=/go/pkg go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux go build -o /github-action-watcher -trimpath -ldflags "-s -w"
FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /github-action-watcher /github-action-watcher
ENTRYPOINT ["/github-action-watcher"]
