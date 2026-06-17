.PHONY: help proto build lint fmt vet test test-unit test-integration test-coverage test-e2e \
        migrate-up migrate-down run-gateway run-user run-all clean

# ========== 默认目标 ==========
.DEFAULT_GOAL := help

# ========== 变量 ==========
GO := go
PROTOC := protoc
DOCKER_COMPOSE := docker compose
GOLANGCI_LINT := golangci-lint

# 服务列表
SERVICES := user-service post-service comment-service mood-service quote-service ai-service notification-service

# ========== 帮助 ==========
help: ## 显示所有可用命令
	@echo "Show Love - 后端构建系统"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

# ========== protobuf ==========
proto: ## 从 proto/ 生成所有 Go gRPC 代码
	@echo "生成 protobuf 代码..."
	@for dir in proto/*/; do \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			$$dir*.proto; \
	done
	@echo "protobuf 代码生成完成"

# ========== 构建 ==========
build: ## 构建所有微服务
	@for svc in $(SERVICES); do \
		echo "构建 $$svc..."; \
		cd services/$$svc && $(GO) build -o ../../bin/$$svc ./cmd/ || exit 1; \
	done
	@echo "构建 gateway..."
	@cd gateway && $(GO) build -o ../bin/gateway ./cmd/
	@echo "全部构建完成"

build-gateway: ## 仅构建网关
	@cd gateway && $(GO) build -o ../bin/gateway ./cmd/

build-user: ## 仅构建用户服务
	@cd services/user-service && $(GO) build -o ../../bin/user-service ./cmd/

# ========== 运行 ==========
run-gateway: ## 启动网关（开发模式）
	@cd gateway && $(GO) run ./cmd/

run-user: ## 启动用户服务
	@cd services/user-service && $(GO) run ./cmd/

run-all: ## 启动所有服务（docker compose up）
	@$(DOCKER_COMPOSE) up -d

# ========== 静态检查 ==========
lint: ## golangci-lint 运行
	@$(GOLANGCI_LINT) run ./...

fmt: ## gofumpt + goimports 格式化所有代码
	@go fmt ./...
	@goimports -w -local showlove .

vet: ## go vet 检查
	@$(GO) vet ./...

# ========== 测试 ==========
test: ## 运行全部后端测试
	@$(GO) test ./... -v -count=1

test-unit: ## 仅单元测试
	@$(GO) test ./... -v -count=1 -short

test-integration: ## 仅集成测试（需 Docker 环境）
	@$(GO) test ./... -v -count=1 -run Integration

test-coverage: ## 带覆盖率报告
	@$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告: coverage.html"

test-e2e: ## 运行 E2E 测试
	@cd test/e2e && $(GO) test ./... -v -count=1

# ========== 数据库 ==========
migrate-up: ## 运行所有服务的数据库迁移
	@echo "运行数据库迁移..."
	@for svc in $(SERVICES); do \
		echo "迁移 $$svc..."; \
		# 各服务独立管理迁移，此处为占位 \
	done

migrate-down: ## 回滚所有迁移
	@echo "回滚数据库迁移..."

# ========== 开发环境 ==========
dev: ## 一键启动开发环境（基础设施 + 全部微服务）
	@bash scripts/dev.sh start

dev-infra: ## 仅启动基础设施（DB/Redis/MinIO/NATS）
	@bash scripts/dev.sh infra

dev-services: ## 仅启动微服务
	@bash scripts/dev.sh services

dev-stop: ## 停止开发环境
	@bash scripts/dev.sh stop

dev-restart: ## 重启开发环境
	@bash scripts/dev.sh restart

dev-status: ## 查看服务状态
	@bash scripts/dev.sh status

dev-health: ## 健康检查
	@bash scripts/dev.sh health

dev-logs: ## 查看网关日志（make dev-logs SVC=user-service）
	@bash scripts/dev.sh logs $(or $(SVC),gateway)

# ========== 清理 ==========
clean: ## 清理构建产物
	@rm -rf bin/ logs/
	@echo "清理完成"
