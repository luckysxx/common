# Common Libraries

`common` 提供多个可复用的 Go 基础模块，供 `user-platform`、`api-gateway`、`go-note` 等服务按需引用。

## Modules

| 模块 | 说明 |
|------|------|
| `common/conf` | 基于 Viper 的基础配置加载与环境隔离 |
| `common/crypto` | 密码哈希与校验工具 |
| `common/errs` | 统一错误码与错误封装 |
| `common/health` | 统一的健康检查接口与探测封装 |
| `common/logger` | 基于 Zap 的日志能力与 Gin 中间件 |
| `common/metrics` | Gin / gRPC 指标采集封装 |
| `common/mq` | 跨服务共享的事件契约、Outbox 抽象与 CDC 集成约定 |
| `common/otel` | OpenTelemetry Tracer 初始化 |
| `common/postgres` | PostgreSQL 数据库连接池初始化与封装 |
| `common/proto` | 共享的 Protobuf 定义与生成代码 |
| `common/ratelimiter` | 多种限流算法实现 |
| `common/redis` | Redis 连接池初始化与统一封装 |
| `common/rpc` | 跨服务 gRPC 客户端封装 |
| `common/trace` | Trace ID 相关工具 |

## Development Notes

- 每个子目录都是独立 Go Module，修改依赖时请在对应模块目录执行 `go mod tidy`。
- `common/proto` 中的 `.pb.go` 文件属于可提交的生成代码，更新 `.proto` 后需要一并提交。
- 本仓库不应提交 IDE 配置、系统缓存文件或本地构建产物。

## Git Hygiene

- 忽略 macOS 缓存文件、编辑器目录和 Go 编译产物。
- 如果误生成了 `.DS_Store` 或临时文件，删除后再提交，避免污染下游服务仓库。
