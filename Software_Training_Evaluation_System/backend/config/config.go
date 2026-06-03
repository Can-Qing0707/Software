package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	LLM      LLMConfig      `yaml:"llm"`
	Upload   UploadConfig   `yaml:"upload"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type LLMConfig struct {
	Provider    string  `yaml:"provider" json:"provider"`
	APIURL      string  `yaml:"api_url" json:"api_url"`
	APIKey      string  `yaml:"api_key" json:"api_key"`
	Model       string  `yaml:"model" json:"model"`
	MaxTokens   int     `yaml:"max_tokens" json:"max_tokens"`
	Temperature float64 `yaml:"temperature" json:"temperature"`
}

type UploadConfig struct {
	Dir     string `yaml:"dir"`
	MaxSize int    `yaml:"max_size"`
}

var AppConfig *Config

func init() {
	AppConfig = &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:    "127.0.0.1",
			Port:    3306,
			User:    "root",
			Charset: "utf8mb4",
		},
		JWT: JWTConfig{
			ExpireHours: 72,
		},
		LLM: LLMConfig{
			Provider:    "openai",
			Model:       "gpt-4o",
			MaxTokens:   4096,
			Temperature: 0.3,
		},
		Upload: UploadConfig{
			Dir:     "./uploads",
			MaxSize: 50,
		},
	}
}

func InitConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	AppConfig = &Config{}
	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		return err
	}
	applyEnvOverrides()
	return nil
}

func applyEnvOverrides() {
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		AppConfig.Database.Password = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		AppConfig.Database.User = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		AppConfig.Database.Host = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		AppConfig.Database.DBName = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		AppConfig.JWT.Secret = v
	}
	if v := os.Getenv("LLM_API_KEY"); v != "" {
		AppConfig.LLM.APIKey = v
	}
	if v := os.Getenv("LLM_API_URL"); v != "" {
		AppConfig.LLM.APIURL = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			AppConfig.Server.Port = p
		}
	}
	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		AppConfig.Upload.Dir = v
	}
}

func (d DatabaseConfig) DSN() string {
	return d.User + ":" + d.Password + "@tcp(" + d.Host + ":" + itoa(d.Port) + ")/" + d.DBName + "?charset=" + d.Charset + "&parseTime=True&loc=Local"
}

func SetLLMConfig(cfg LLMConfig) {
	AppConfig.LLM = cfg
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
