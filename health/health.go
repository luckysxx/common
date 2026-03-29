// Package health 提供标准化的健康检查端点。
//
// 设计遵循 Kubernetes 探针规范：
//   - /healthz (Liveness)  — 进程存活即返回 200，用于检测死锁或不可恢复的状态
//   - /readyz  (Readiness) — 所有依赖就绪才返回 200，用于判断是否可以接收流量
//
// 用法：
//
//	checker := health.NewChecker()
//	checker.AddCheck("postgres", func(ctx context.Context) error { return db.PingContext(ctx) })
//	checker.AddCheck("redis", func(ctx context.Context) error { return rdb.Ping(ctx).Err() })
//	checker.Register(r) // r 是 *gin.Engine
package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CheckFunc 是单个依赖的健康检查函数，返回 nil 表示健康。
type CheckFunc func(ctx context.Context) error

// Checker 管理一组依赖检查项，提供 /healthz 和 /readyz 端点。
type Checker struct {
	mu     sync.RWMutex
	checks map[string]CheckFunc
}

// NewChecker 创建一个新的健康检查器。
func NewChecker() *Checker {
	return &Checker{
		checks: make(map[string]CheckFunc),
	}
}

// AddCheck 注册一个命名的依赖检查项。
// name 用于在响应中标识具体哪个依赖不健康。
func (c *Checker) AddCheck(name string, fn CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = fn
}

// Register 将 /healthz 和 /readyz 路由注册到 Gin 引擎上。
// 这两个端点注册在 metrics 和业务中间件之前，避免被限流或鉴权拦截。
func (c *Checker) Register(r *gin.Engine) {
	r.GET("/healthz", c.liveness)
	r.GET("/readyz", c.readiness)
}

// liveness 存活探针：进程在跑就返回 200。
// K8s 用这个判断是否需要重启容器。
func (c *Checker) liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// readiness 就绪探针：逐一检查所有依赖，全部通过才返回 200。
// K8s 用这个判断是否把流量路由到该 Pod。
func (c *Checker) readiness(ctx *gin.Context) {
	c.mu.RLock()
	checks := make(map[string]CheckFunc, len(c.checks))
	for k, v := range c.checks {
		checks[k] = v
	}
	c.mu.RUnlock()

	// 给每个检查项 2 秒超时，防止慢依赖拖垮探针
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	results := make(map[string]string, len(checks))
	allHealthy := true

	for name, fn := range checks {
		if err := fn(checkCtx); err != nil {
			results[name] = err.Error()
			allHealthy = false
		} else {
			results[name] = "ok"
		}
	}

	status := http.StatusOK
	overall := "ready"
	if !allHealthy {
		status = http.StatusServiceUnavailable
		overall = "not_ready"
	}

	ctx.JSON(status, gin.H{
		"status": overall,
		"checks": results,
	})
}
