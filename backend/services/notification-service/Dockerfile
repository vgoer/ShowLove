FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/
COPY pkg/events/ ./pkg/events/
WORKDIR /app/services/notification-service
RUN go mod download
WORKDIR /app
COPY . .
WORKDIR /app/services/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd/
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /notification-service /usr/local/bin/notification-service
EXPOSE 50057
CMD ["notification-service"]
