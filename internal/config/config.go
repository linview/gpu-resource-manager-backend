package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	NATS     NATSConfig
	Tinkerbell TinkerbellConfig
	K8s      K8sConfig
	LogLevel string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NATSConfig NATS配置
type NATSConfig struct {
	URL string
}

// TinkerbellConfig Tinkerbell配置
type TinkerbellConfig struct {
	URL      string
	Username string
	Password string
}

// K8sConfig Kubernetes配置
type K8sConfig struct {
	ConfigPath string
	Namespace  string
}

// Load 加载配置
func Load() *Config {
	// 加载.env文件
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "gpu_management"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		Tinkerbell: TinkerbellConfig{
			URL:      getEnv("TINKERBELL_URL", "http://localhost:50061"),
			Username: getEnv("TINKERBELL_USERNAME", "admin"),
			Password: getEnv("TINKERBELL_PASSWORD", "password"),
		},
		K8s: K8sConfig{
			ConfigPath: getEnv("K8S_CONFIG_PATH", ""),
			Namespace:  getEnv("K8S_NAMESPACE", "default"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
