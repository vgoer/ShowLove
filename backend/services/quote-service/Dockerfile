FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY services/quote-service/go.mod services/quote-service/go.sum ./services/quote-service/
WORKDIR /app/services/quote-service
RUN go mod download
WORKDIR /app
COPY . .
WORKDIR /app/services/quote-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /quote-service ./cmd/
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /quote-service /usr/local/bin/quote-service
EXPOSE 50055
CMD ["quote-service"]
