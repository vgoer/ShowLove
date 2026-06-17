# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY services/post-service/go.mod services/post-service/go.sum ./services/post-service/
COPY pkg/errcode/ ./pkg/errcode/
COPY pkg/events/ ./pkg/events/

WORKDIR /app/services/post-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/services/post-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /post-service ./cmd/

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /post-service /usr/local/bin/post-service
EXPOSE 50052
CMD ["post-service"]
