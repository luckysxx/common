// Package conf 提供统一的 Viper 配置加载函数和少量通用配置类型。
//
// 设计原则：
//   - 基础设施配置类型（Redis / Postgres）定义在各自的 common 包中
//   - 只保留无对应包的通用类型（OTel / Server / IDGenerator）
//   - 各服务通过嵌入公共类型来组合自己的 Config 结构体
package conf

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

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
