# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.work and module files for dependency resolution
COPY go.work go.work.sum ./
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/
COPY pkg/jwt/go.mod pkg/jwt/go.sum ./pkg/jwt/

# Download dependencies
WORKDIR /app/services/user-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY pkg/ ./pkg/
COPY services/user-service/ ./services/user-service/

# Build
WORKDIR /app/services/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /user-service /usr/local/bin/user-service

EXPOSE 50051

CMD ["user-service"]
