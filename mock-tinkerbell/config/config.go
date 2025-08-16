package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Hardware HardwareConfig `mapstructure:"hardware"`
	Workflow WorkflowConfig `mapstructure:"workflow"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type HardwareConfig struct {
	DefaultGPUSpecs []GPUSpec `mapstructure:"default_gpu_specs"`
	MockDelay       int       `mapstructure:"mock_delay_ms"`
}

type GPUSpec struct {
	Model        string `mapstructure:"model"`
	MemoryGB     int    `mapstructure:"memory_gb"`
	DriverVersion string `mapstructure:"driver_version"`
	CUDAVersion  string `mapstructure:"cuda_version"`
	Count        int    `mapstructure:"count"`
}

type WorkflowConfig struct {
	DefaultTimeout int `mapstructure:"default_timeout_seconds"`
	StepDelay      int `mapstructure:"step_delay_ms"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file, using defaults: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("hardware.mock_delay_ms", 1000)
	viper.SetDefault("workflow.default_timeout_seconds", 300)
	viper.SetDefault("workflow.step_delay_ms", 2000)

	// 默认GPU规格
	viper.SetDefault("hardware.default_gpu_specs", []map[string]interface{}{
		{
			"model":          "NVIDIA A100",
			"memory_gb":      80,
			"driver_version": "525.85.05",
			"cuda_version":   "12.0",
			"count":          4,
		},
		{
			"model":          "NVIDIA V100",
			"memory_gb":      32,
			"driver_version": "525.85.05",
			"driver_version": "11.8",
			"count":          8,
		},
		{
			"model":          "NVIDIA RTX 4090",
			"memory_gb":      24,
			"driver_version": "525.85.05",
			"cuda_version":   "12.0",
			"count":          2,
		},
	})
}
