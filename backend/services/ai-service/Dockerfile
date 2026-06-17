FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY services/ai-service/go.mod services/ai-service/go.sum ./services/ai-service/
COPY pkg/events/ ./pkg/events/
WORKDIR /app/services/ai-service
RUN go mod download
WORKDIR /app
COPY . .
WORKDIR /app/services/ai-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai-service ./cmd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /ai-service /usr/local/bin/ai-service
EXPOSE 50056
CMD ["ai-service"]
