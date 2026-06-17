FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY services/mood-service/go.mod services/mood-service/go.sum ./services/mood-service/
WORKDIR /app/services/mood-service
RUN go mod download
WORKDIR /app
COPY . .
WORKDIR /app/services/mood-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /mood-service ./cmd/
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /mood-service /usr/local/bin/mood-service
EXPOSE 50054
CMD ["mood-service"]
