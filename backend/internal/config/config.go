package config

import "os"

// Role of this struct = only runtime settings container
type Config struct {
	Port string
	DB   string

	JWTSecret string

	Redis string
	Minio string

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

	return Config{
		Port: port,
		DB:   os.Getenv("DATABASE_URL"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		Redis: os.Getenv("REDIS_URL"),
		Minio: os.Getenv("MINIO_URL"),

		BaseURL:      os.Getenv("KEYCLOAK_BASE_URL"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),

		FrontendURL:    os.Getenv("FRONTEND_URL"),
		BackendBaseURL: os.Getenv("BACKEND_URL"),
	}
}
