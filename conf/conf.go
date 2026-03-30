// Package conf 提供微服务间共享的通用配置类型定义和统一的 Viper 配置加载函数。
//
// 设计原则：
//   - 只抽取各服务完全重复的配置类型（Database / Redis / OTel / Server / IDGenerator）
//   - 服务专有的配置（如 JWTConfig / KafkaConfig / ChatConfig）仍由各服务自行定义
//   - 各服务通过嵌入公共类型来组合自己的 Config 结构体
package conf

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// DatabaseConfig 通用数据库配置
type DatabaseConfig struct {
	Driver      string `mapstructure:"driver"`
	Source      string `mapstructure:"source"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

// RedisConfig 通用 Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// OTelConfig 通用可观测性配置
type OTelConfig struct {
	JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
	ServiceName    string `mapstructure:"service_name"`
}

// ServerConfig 通用 HTTP 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// IDGeneratorConfig 分布式 ID 生成器配置
type IDGeneratorConfig struct {
	Addr string `mapstructure:"addr"`
}

// Load 加载配置到目标结构体。
//
// 统一了各服务重复的加载逻辑：godotenv → viper 读 YAML → 环境变量覆盖 → Unmarshal。
// target 必须是指针类型（例如 &Config{}），其字段可以嵌入上面的公共类型，
// 也可以包含服务专有的类型。
func Load(target any) {
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config.yaml found, relying entirely on ENV variables: %v", err)
	}

	if err := viper.Unmarshal(target); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}
