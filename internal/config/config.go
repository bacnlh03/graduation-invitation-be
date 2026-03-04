package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort           string   `mapstructure:"APP_PORT"`
	DBHost            string   `mapstructure:"DB_HOST"`
	DBPort            string   `mapstructure:"DB_PORT"`
	DBUser            string   `mapstructure:"DB_USER"`
	DBPassword        string   `mapstructure:"DB_PASSWORD"`
	DBName            string   `mapstructure:"DB_NAME"`
	DBSSLMode         string   `mapstructure:"DB_SSLMODE"`
	AdminUsername     string   `mapstructure:"ADMIN_USERNAME"`
	AdminPasswordHash string   `mapstructure:"ADMIN_PASSWORD_HASH"`
	JWTSecret         string   `mapstructure:"JWT_SECRET"`
	AppEnv            string   `mapstructure:"APP_ENV"`
	CORSAllowed       []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	UploadDir         string   `mapstructure:"UPLOAD_DIR"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Không tìm thấy file .env, sẽ dùng biến môi trường hệ thống")
	}

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("UPLOAD_DIR", "uploads")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Không thể map cấu hình:", err)
	}

	return &config
}
