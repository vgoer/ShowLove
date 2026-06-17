# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY pkg/jwt/go.mod pkg/jwt/go.sum ./pkg/jwt/
COPY pkg/validator/go.mod pkg/validator/go.sum ./pkg/validator/
COPY pkg/errcode/go.mod pkg/errcode/go.sum ./pkg/errcode/

WORKDIR /app/gateway
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./cmd/

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /gateway /usr/local/bin/gateway

EXPOSE 8080

CMD ["gateway"]
