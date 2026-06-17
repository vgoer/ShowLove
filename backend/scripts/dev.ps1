<#
.SYNOPSIS
    Show Love 本地开发环境管理脚本 (PowerShell)
.DESCRIPTION
    一键启动/停止所有后端基础设施和微服务。
.PARAMETER Action
    start    - 启动基础设施 + 全部微服务
    stop     - 停止全部服务
    restart  - 重启全部服务
    status   - 查看服务运行状态
    infra    - 仅启动基础设施 (DB/Redis/MinIO/NATS)
    services - 仅启动微服务 (需要基础设施已运行)
    logs     - 查看指定服务的日志
.EXAMPLE
    .\scripts\dev.ps1 start          # 一键启动全部
    .\scripts\dev.ps1 stop           # 停止全部
    .\scripts\dev.ps1 status         # 查看状态
    .\scripts\dev.ps1 logs gateway   # 查看网关日志
#>

param(
    [Parameter(Position = 0)]
    [ValidateSet("start", "stop", "restart", "status", "infra", "services", "logs", "health")]
    [string]$Action = "start",

    [Parameter(Position = 1)]
    [string]$ServiceName = ""
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Split-Path -Parent $ScriptDir

# 服务端口映射
$Services = @{
    "postgres"    = @{ Port = 5432; HealthUrl = "";   Type = "infra" }
    "redis"       = @{ Port = 6379; HealthUrl = "";   Type = "infra" }
    "minio"       = @{ Port = 9000; HealthUrl = "http://localhost:9000/minio/health/live"; Type = "infra" }
    "nats"        = @{ Port = 4222; HealthUrl = "http://localhost:8222/healthz"; Type = "infra" }
    "gateway"     = @{ Port = 8080; HealthUrl = "http://localhost:8080/api/v1/health"; Type = "app" }
    "user"        = @{ Port = 50051; HealthUrl = ""; Type = "app"; Dir = "services/user-service" }
    "post"        = @{ Port = 50052; HealthUrl = ""; Type = "app"; Dir = "services/post-service" }
    "comment"     = @{ Port = 50053; HealthUrl = ""; Type = "app"; Dir = "services/comment-service" }
    "mood"        = @{ Port = 50054; HealthUrl = ""; Type = "app"; Dir = "services/mood-service" }
    "quote"       = @{ Port = 50055; HealthUrl = ""; Type = "app"; Dir = "services/quote-service" }
    "ai"          = @{ Port = 50056; HealthUrl = ""; Type = "app"; Dir = "services/ai-service" }
    "notification" = @{ Port = 50057; HealthUrl = ""; Type = "app"; Dir = "services/notification-service" }
}

# 颜色输出
function Write-Info  { Write-Host "[INFO] " -NoNewline -ForegroundColor Cyan; Write-Host $args }
function Write-OK    { Write-Host "[OK]   " -NoNewline -ForegroundColor Green; Write-Host $args }
function Write-Warn  { Write-Host "[WARN] " -NoNewline -ForegroundColor Yellow; Write-Host $args }
function Write-ErrorMsg { Write-Host "[ERROR]" -NoNewline -ForegroundColor Red; Write-Host $args }

# 确保在 backend 目录
Push-Location $BackendDir

# ========== 端口检查 ==========
function Test-PortInUse {
    param([int]$Port)
    $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    return ($null -ne $connections -and $connections.Count -gt 0)
}

# ========== Docker 基础设施 ==========
function Start-Infra {
    Write-Info "启动 Docker 基础设施..."

    # 检查 Docker 是否运行
    $dockerRunning = docker info 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-ErrorMsg "Docker 未运行，请先启动 Docker Desktop"
        return $false
    }

    docker compose up -d postgres redis minio minio-init nats 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-ErrorMsg "Docker 启动失败"
        return $false
    }

    Write-Info "等待基础设施就绪..."
    $timeout = 60
    $elapsed = 0
    $services = @("postgres", "redis", "minio", "nats")

    while ($elapsed -lt $timeout) {
        $allHealthy = $true
        foreach ($svc in $services) {
            $status = docker compose ps --format json $svc 2>$null | ConvertFrom-Json -ErrorAction SilentlyContinue
            if (-not $status -or $status.Health -ne "healthy") {
                $allHealthy = $false
                break
            }
        }
        if ($allHealthy) { break }
        Start-Sleep -Seconds 2
        $elapsed += 2
        Write-Host "." -NoNewline
    }
    Write-Host ""

    if ($elapsed -ge $timeout) {
        Write-Warn "部分基础设施可能未就绪 (超时 ${timeout}s)"
    } else {
        Write-OK "基础设施全部就绪 (${elapsed}s)"
    }

    # 初始化数据库（如果 init-multi-db.sh 未自动执行）
    Write-Info "检查数据库..."
    $pgUser = $env:POSTGRES_USER
    if (-not $pgUser) { $pgUser = "user" }
    $dbs = @("users_db", "posts_db", "comments_db", "moods_db", "quotes_db", "notifications_db")
    foreach ($db in $dbs) {
        $result = docker exec showlove-postgres psql -U $pgUser -lqt 2>$null | Select-String $db
        if (-not $result) {
            Write-Info "创建数据库: $db"
            docker exec showlove-postgres psql -U $pgUser -c "CREATE DATABASE $db" 2>$null | Out-Null
        }
    }
    Write-OK "数据库检查完成"

    return $true
}

function Stop-Infra {
    Write-Info "停止 Docker 基础设施..."
    docker compose stop postgres redis minio nats 2>&1 | Out-Null
    Write-OK "基础设施已停止"
}

# ========== 微服务管理 ==========
$RunningServices = @{}

function Start-AppService {
    param([string]$Name, [string]$Dir)

    $port = $Services[$Name].Port
    if (Test-PortInUse $port) {
        Write-Warn "$Name 端口 $port 已被占用，跳过"
        return
    }

    Write-Info "启动 $Name (:$port)..."

    $logFile = Join-Path $BackendDir "logs" "$Name.log"
    New-Item -ItemType Directory -Force -Path (Split-Path $logFile) | Out-Null

    $proc = Start-Process -FilePath "go" `
        -ArgumentList "run ./cmd/" `
        -WorkingDirectory (Join-Path $BackendDir $Dir) `
        -RedirectStandardOutput $logFile `
        -RedirectStandardError $logFile `
        -NoNewWindow `
        -PassThru

    $RunningServices[$Name] = $proc
    Write-OK "$Name 启动中 (PID: $($proc.Id), 日志: $logFile)"
    Start-Sleep -Seconds 1
}

function Start-AllServices {
    Write-Info "启动全部微服务 + 网关..."

    # 按依赖顺序启动
    $startOrder = @("user", "post", "comment", "mood", "quote", "ai", "notification", "gateway")

    foreach ($svc in $startOrder) {
        if ($svc -eq "gateway") {
            $proc = Start-Process -FilePath "go" `
                -ArgumentList "run ./gateway/cmd/" `
                -WorkingDirectory $BackendDir `
                -RedirectStandardOutput "$BackendDir\logs\gateway.log" `
                -RedirectStandardError "$BackendDir\logs\gateway.log" `
                -NoNewWindow `
                -PassThru
            $RunningServices["gateway"] = $proc
            Write-OK "gateway 启动中 (PID: $($proc.Id))"
        } else {
            $dir = $Services[$svc].Dir
            Start-AppService -Name $svc -Dir $dir
        }
    }

    Write-Info "等待网关就绪..."
    $timeout = 30
    $elapsed = 0
    while ($elapsed -lt $timeout) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/health" -TimeoutSec 2 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-OK "网关就绪！http://localhost:8080"
                return
            }
        } catch {}
        Start-Sleep -Seconds 2
        $elapsed += 2
        Write-Host "." -NoNewline
    }
    Write-Host ""
    Write-Warn "网关可能尚未就绪，查看日志: logs/gateway.log"
}

function Stop-AllServices {
    Write-Info "停止全部微服务..."

    # 通过端口杀进程
    $appPorts = @(8080, 50051, 50052, 50053, 50054, 50055, 50056, 50057)
    foreach ($port in $appPorts) {
        $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        if ($conn) {
            $procId = $conn.OwningProcess
            try {
                Stop-Process -Id $procId -Force -ErrorAction SilentlyContinue
                Write-OK "已停止端口 $port 进程 (PID: $procId)"
            } catch {
                Write-Warn "无法停止端口 $port 进程"
            }
        }
    }

    foreach ($key in $RunningServices.Keys) {
        $proc = $RunningServices[$key]
        if (-not $proc.HasExited) {
            $proc.Kill()
            Write-OK "已停止 $key (PID: $($proc.Id))"
        }
    }
    $RunningServices.Clear()
}

# ========== 状态检查 ==========
function Show-Status {
    Write-Host ""
    Write-Host "═══ Show Love 服务状态 ═══" -ForegroundColor Cyan
    Write-Host ""

    # Docker 服务
    Write-Host "▸ 基础设施 (Docker)" -ForegroundColor Yellow
    $infraSvcs = @("postgres", "redis", "minio", "nats")
    foreach ($svc in $infraSvcs) {
        $status = docker compose ps --format json $svc 2>$null | ConvertFrom-Json -ErrorAction SilentlyContinue
        if ($status) {
            $health = $status.Health
            $color = if ($health -eq "healthy") { "Green" } else { "Yellow" }
            Write-Host "  $svc" -NoNewline
            Write-Host " [$health]" -ForegroundColor $color
        } else {
            Write-Host "  $svc" -NoNewline
            Write-Host " [stopped]" -ForegroundColor Red
        }
    }

    # 微服务
    Write-Host ""
    Write-Host "▸ 微服务" -ForegroundColor Yellow
    $appSvcs = @("user", "post", "comment", "mood", "quote", "ai", "notification", "gateway")
    foreach ($svc in $appSvcs) {
        $port = $Services[$svc].Port
        if (Test-PortInUse $port) {
            Write-Host "  $svc (:$port)" -NoNewline
            Write-Host " [running]" -ForegroundColor Green
        } else {
            Write-Host "  $svc (:$port)" -NoNewline
            Write-Host " [stopped]" -ForegroundColor Red
        }
    }

    # Gateway health
    Write-Host ""
    Write-Host "▸ API 网关健康检查" -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health" -TimeoutSec 3 -ErrorAction Stop
        Write-Host "  GET /health → " -NoNewline
        Write-Host "200 OK" -ForegroundColor Green
        Write-Host "  $($response | ConvertTo-Json -Compress)"
    } catch {
        Write-Host "  GET /health → " -NoNewline
        Write-Host "不可达" -ForegroundColor Red
    }

    Write-Host ""
    Write-Host "═══ 快捷命令 ═══" -ForegroundColor Cyan
    Write-Host "  .\scripts\dev.ps1 start    → 一键启动全部"
    Write-Host "  .\scripts\dev.ps1 stop     → 停止全部"
    Write-Host "  .\scripts\dev.ps1 logs <name> → 查看日志"
    Write-Host "  curl http://localhost:8080/api/v1/health"
    Write-Host ""
}

# ========== 日志查看 ==========
function Show-Logs {
    param([string]$Name)

    if (-not $Name) {
        Write-Info "可用服务: gateway, user, post, comment, mood, quote, ai, notification"
        Write-Info "用法: .\scripts\dev.ps1 logs gateway"
        return
    }

    $logFile = Join-Path $BackendDir "logs" "$Name.log"
    if (Test-Path $logFile) {
        Write-Info "查看 $Name 日志 (Ctrl+C 退出)..."
        Get-Content $logFile -Wait -Tail 50
    } elseif ($Services[$Name].Type -eq "infra") {
        Write-Info "查看 Docker 日志..."
        docker compose logs -f --tail 50 $Name
    } else {
        Write-Warn "日志文件不存在: $logFile"
    }
}

# ========== 健康检查 ==========
function Test-Health {
    Write-Info "检查所有服务健康状态..."

    # Docker
    Write-Info "Docker 服务:"
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}" 2>$null

    # Gateway
    Write-Info "API 网关:"
    try {
        $r = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health" -TimeoutSec 3
        Write-OK "网关正常: $($r | ConvertTo-Json -Compress)"
    } catch {
        Write-ErrorMsg "网关不可达"
    }

    # 端口列表
    Write-Info "端口监听:"
    $ports = @(5432, 6379, 9000, 4222, 8080, 50051, 50052, 50053, 50054, 50055, 50056, 50057)
    foreach ($p in $ports) {
        if (Test-PortInUse $p) {
            Write-OK ":$p 已监听"
        } else {
            Write-Warn ":$p 未监听"
        }
    }
}

# ========== 主流程 ==========
switch ($Action) {
    "start" {
        Write-Host ""
        Write-Host "🌸 显出爱心 - 开发环境启动" -ForegroundColor Magenta
        Write-Host ""

        if (Start-Infra) {
            Start-AllServices
        }

        Write-Host ""
        Write-Host "═══ 启动完成 ═══" -ForegroundColor Green
        Write-Host "  API 网关: http://localhost:8080"
        Write-Host "  MinIO 控制台: http://localhost:9001 (minioadmin/minioadmin)"
        Write-Host "  NATS 监控: http://localhost:8222"
        Write-Host "  日志目录: backend/logs/"
        Write-Host ""
        Write-Host "  测试命令:"
        Write-Host "    curl http://localhost:8080/api/v1/health"
        Write-Host "    .\scripts\dev.ps1 status"
        Write-Host "    .\scripts\dev.ps1 logs gateway"
        Write-Host ""
    }
    "stop" {
        Write-Info "停止 Show Love 开发环境..."
        Stop-AllServices
        Stop-Infra
        Write-OK "全部已停止"
    }
    "restart" {
        Write-Info "重启 Show Love 开发环境..."
        Stop-AllServices
        Start-Sleep 2
        if (Start-Infra) {
            Start-AllServices
        }
        Write-OK "重启完成"
    }
    "infra" {
        Write-Host ""
        Write-Host "🌸 启动基础设施" -ForegroundColor Magenta
        Write-Host ""
        Start-Infra | Out-Null
        Write-Host ""
        Write-Host "基础设施就绪: PostgreSQL(5432) Redis(6379) MinIO(9000) NATS(4222)" -ForegroundColor Green
        Write-Host ""
    }
    "services" {
        Write-Host ""
        Write-Host "🌸 启动微服务" -ForegroundColor Magenta
        Write-Host ""
        Start-AllServices
    }
    "status" {
        Show-Status
    }
    "logs" {
        Show-Logs -Name $ServiceName
    }
    "health" {
        Test-Health
    }
}

Pop-Location
