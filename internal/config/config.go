package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	// 服务器配置
	Port string

	// 数据库配置
	DBPath string

	// LLM配置
	LLMProvider string
	LLMAPIKey   string
	LLMAPIURL   string
	LLMModel    string

	// JWT配置
	JWTSecret string

	// 日志配置
	LogLevel string
}

// Load 加载配置
func Load() *Config {
	// 加载环境变量文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./qa_database.db"),
		LLMProvider: getEnv("LLM_PROVIDER", "openai"),
		LLMAPIKey:   getEnv("LLM_API_KEY", ""),
		LLMAPIURL:   getEnv("LLM_API_URL", ""),
		LLMModel:    getEnv("LLM_MODEL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetLLMMode 获取LLM运行模式
func (c *Config) GetLLMMode() string {
	if c.LLMAPIKey != "" {
		return "真实API模式"
	}
	return "模拟演示模式"
}
