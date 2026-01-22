package config

import (
	"os"
	"strconv"
)

// Role of this struct = only runtime settings container
type Config struct {
	Port string
	DB   string

	JWTSecret string

	Redis          string
	Minio          string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	MinioPublicURL string

	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string

	FrontendURL    string
	BackendBaseURL string
}

func Load() Config {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	minioUseSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	return Config{
		Port: port,
		DB:   os.Getenv("DATABASE_URL"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		Redis:          os.Getenv("REDIS_URL"),
		Minio:          os.Getenv("MINIO_URL"),
		MinioEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioBucket:    os.Getenv("MINIO_BUCKET"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioPublicURL: os.Getenv("MINIO_PUBLIC_URL"),
		MinioUseSSL:    minioUseSSL,

		BaseURL:      os.Getenv("KEYCLOAK_BASE_URL"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),

		FrontendURL:    os.Getenv("FRONTEND_URL"),
		BackendBaseURL: os.Getenv("BACKEND_URL"),
	}
}
