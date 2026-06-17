#!/usr/bin/env bash
# =============================================================================
# Show Love 本地开发环境管理脚本 (Bash / Git Bash / Linux / Mac)
#
# 用法:
#   ./scripts/dev.sh start         一键启动全部
#   ./scripts/dev.sh stop          停止全部
#   ./scripts/dev.sh restart       重启全部
#   ./scripts/dev.sh status        查看状态
#   ./scripts/dev.sh infra         仅启动基础设施
#   ./scripts/dev.sh services      仅启动微服务
#   ./scripts/dev.sh logs <name>   查看服务日志
#   ./scripts/dev.sh health        健康检查
# =============================================================================

set -e

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG_DIR="$BACKEND_DIR/logs"

ACTION="${1:-start}"
SERVICE_NAME="${2:-}"

info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()    { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }

# 端口→进程清理
kill_port() {
    local port=$1
    local pid
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        # Git Bash on Windows
        pid=$(netstat -ano 2>/dev/null | grep ":$port " | grep LISTENING | awk '{print $5}' | head -1)
    else
        pid=$(lsof -ti :"$port" 2>/dev/null)
    fi
    if [[ -n "$pid" ]]; then
        kill "$pid" 2>/dev/null && ok "已停止端口 $port (PID: $pid)"
    fi
}

# ========== 基础设施 ==========
start_infra() {
    info "启动 Docker 基础设施..."
    if ! docker info &>/dev/null; then
        error "Docker 未运行，请先启动 Docker Desktop"
        return 1
    fi

    docker compose up -d postgres redis minio minio-init nats 2>&1 | tail -5

    info "等待基础设施就绪..."
    local timeout=60 elapsed=0
    while ((elapsed < timeout)); do
        local all_healthy=true
        for svc in postgres redis minio nats; do
            local health
            health=$(docker compose ps --format json "$svc" 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('Health',''))" 2>/dev/null || echo "")
            if [[ "$health" != "healthy" ]]; then
                all_healthy=false
                break
            fi
        done
        if $all_healthy; then break; fi
        sleep 2
        ((elapsed+=2))
        printf "."
    done
    echo ""

    if ((elapsed >= timeout)); then
        warn "部分基础设施可能未就绪 (超时 ${timeout}s)"
    else
        ok "基础设施全部就绪 (${elapsed}s)"
    fi

    # 确保数据库存在
    info "检查数据库..."
    local dbs=("users_db" "posts_db" "comments_db" "moods_db" "quotes_db" "notifications_db")
    for db in "${dbs[@]}"; do
        if ! docker exec showlove-postgres psql -U "${POSTGRES_USER:-user}" -lqt 2>/dev/null | grep -q "$db"; then
            info "创建数据库: $db"
            docker exec showlove-postgres psql -U "${POSTGRES_USER:-user}" -c "CREATE DATABASE $db" 2>/dev/null || true
        fi
    done
    ok "数据库检查完成"
    return 0
}

stop_infra() {
    info "停止 Docker 基础设施..."
    docker compose stop postgres redis minio nats 2>&1 | tail -3
    ok "基础设施已停止"
}

# ========== 微服务 ==========
start_services() {
    info "启动全部微服务 + 网关..."
    mkdir -p "$LOG_DIR"

    # 按依赖顺序启动
    local services=(
        "user-service:50051:services/user-service"
        "post-service:50052:services/post-service"
        "comment-service:50053:services/comment-service"
        "mood-service:50054:services/mood-service"
        "quote-service:50055:services/quote-service"
        "ai-service:50056:services/ai-service"
        "notification-service:50057:services/notification-service"
    )

    for svc_def in "${services[@]}"; do
        IFS=':' read -r name port dir <<< "$svc_def"
        if [[ $(lsof -ti :"$port" 2>/dev/null || netstat -ano 2>/dev/null | grep -c ":$port ") -gt 0 ]]; then
            warn "$name 端口 $port 已被占用，跳过"
            continue
        fi
        info "启动 $name (:$port)..."
        (cd "$BACKEND_DIR/$dir" && go run ./cmd/ >> "$LOG_DIR/$name.log" 2>&1) &
        ok "$name 启动中 (日志: $LOG_DIR/$name.log)"
        sleep 1
    done

    # Gateway
    info "启动 gateway (:8080)..."
    (cd "$BACKEND_DIR" && go run ./gateway/cmd/ >> "$LOG_DIR/gateway.log" 2>&1) &
    ok "gateway 启动中 (日志: $LOG_DIR/gateway.log)"

    # 等待网关就绪
    info "等待网关就绪..."
    local timeout=30 elapsed=0
    while ((elapsed < timeout)); do
        if curl -s http://localhost:8080/api/v1/health &>/dev/null; then
            ok "网关就绪！http://localhost:8080"
            return 0
        fi
        sleep 2
        ((elapsed+=2))
        printf "."
    done
    echo ""
    warn "网关可能尚未就绪，查看日志: $LOG_DIR/gateway.log"
}

stop_services() {
    info "停止全部微服务..."
    local ports=(8080 50051 50052 50053 50054 50055 50056 50057)
    for port in "${ports[@]}"; do
        kill_port "$port"
    done
    ok "微服务已停止"
}

# ========== 状态 ==========
show_status() {
    echo ""
    echo -e "${MAGENTA}═══ Show Love 服务状态 ═══${NC}"
    echo ""

    echo -e "${YELLOW}▸ 基础设施 (Docker)${NC}"
    for svc in postgres redis minio nats; do
        local status
        status=$(docker compose ps --format json "$svc" 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('Health','stopped'))" 2>/dev/null || echo "stopped")
        case "$status" in
            healthy) echo -e "  $svc ${GREEN}[$status]${NC}" ;;
            stopped) echo -e "  $svc ${RED}[$status]${NC}" ;;
            *)       echo -e "  $svc ${YELLOW}[$status]${NC}" ;;
        esac
    done

    echo ""
    echo -e "${YELLOW}▸ 微服务${NC}"
    local apps=("user:50051" "post:50052" "comment:50053" "mood:50054" "quote:50055" "ai:50056" "notification:50057" "gateway:8080")
    for app_def in "${apps[@]}"; do
        IFS=':' read -r name port <<< "$app_def"
        if lsof -ti :"$port" &>/dev/null 2>&1 || netstat -ano 2>/dev/null | grep -q ":$port "; then
            echo -e "  $name (:$port) ${GREEN}[running]${NC}"
        else
            echo -e "  $name (:$port) ${RED}[stopped]${NC}"
        fi
    done

    echo ""
    echo -e "${YELLOW}▸ API 网关健康检查${NC}"
    if curl -s http://localhost:8080/api/v1/health &>/dev/null; then
        echo -e "  GET /health → ${GREEN}200 OK${NC}"
        curl -s http://localhost:8080/api/v1/health | python3 -m json.tool 2>/dev/null || true
    else
        echo -e "  GET /health → ${RED}不可达${NC}"
    fi
    echo ""
}

# ========== 主流程 ==========
cd "$BACKEND_DIR"

case "$ACTION" in
    start)
        echo ""
        echo -e "${MAGENTA}🌸 显出爱心 - 开发环境启动${NC}"
        echo ""
        start_infra && start_services
        echo ""
        echo -e "${GREEN}═══ 启动完成 ═══${NC}"
        echo "  API 网关: http://localhost:8080"
        echo "  MinIO 控制台: http://localhost:9001 (minioadmin/minioadmin)"
        echo "  NATS 监控: http://localhost:8222"
        echo "  日志目录: $LOG_DIR"
        echo ""
        ;;
    stop)
        info "停止 Show Love 开发环境..."
        stop_services
        stop_infra
        ok "全部已停止"
        ;;
    restart)
        info "重启 Show Love 开发环境..."
        stop_services
        sleep 2
        start_infra && start_services
        ok "重启完成"
        ;;
    infra)
        echo ""
        echo -e "${MAGENTA}🌸 启动基础设施${NC}"
        echo ""
        start_infra
        echo ""
        ok "基础设施就绪: PostgreSQL(5432) Redis(6379) MinIO(9000) NATS(4222)"
        ;;
    services)
        echo ""
        echo -e "${MAGENTA}🌸 启动微服务${NC}"
        echo ""
        start_services
        ;;
    status)
        show_status
        ;;
    logs)
        if [[ -z "$SERVICE_NAME" ]]; then
            info "可用服务: gateway, user-service, post-service, comment-service, mood-service, quote-service, ai-service, notification-service"
            info "用法: ./scripts/dev.sh logs gateway"
        elif [[ -f "$LOG_DIR/$SERVICE_NAME.log" ]]; then
            info "查看 $SERVICE_NAME 日志 (Ctrl+C 退出)..."
            tail -f "$LOG_DIR/$SERVICE_NAME.log"
        else
            warn "日志文件不存在: $LOG_DIR/$SERVICE_NAME.log"
        fi
        ;;
    health)
        info "健康检查..."
        docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}" 2>/dev/null
        echo ""
        if curl -s http://localhost:8080/api/v1/health; then
            echo ""
            ok "网关正常"
        else
            error "网关不可达"
        fi
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status|infra|services|logs|health} [service-name]"
        exit 1
        ;;
esac
