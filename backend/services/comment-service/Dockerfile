FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY services/comment-service/go.mod services/comment-service/go.sum ./services/comment-service/
WORKDIR /app/services/comment-service
RUN go mod download
WORKDIR /app
COPY . .
WORKDIR /app/services/comment-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /comment-service ./cmd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /comment-service /usr/local/bin/comment-service
EXPOSE 50053
CMD ["comment-service"]
