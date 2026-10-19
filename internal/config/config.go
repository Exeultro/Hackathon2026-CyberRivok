package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env       string
	HTTP      HTTPConfig
	DB        DBConfig
	JWT       JWTConfig
	Argon2    Argon2Config
	CORS      CORSConfig
	StaticDir string
}

type HTTPConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	DSN      string
	MaxConns int32
	MinConns int32
}

type JWTConfig struct {
	Secret     []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

type Argon2Config struct {
	MemoryKiB   uint32
	Iterations  uint32
	Parallelism uint8
	SaltLen     uint32
	KeyLen      uint32
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env: getenv("APP_ENV", "dev"),
		HTTP: HTTPConfig{
			Addr:         getenv("HTTP_ADDR", ":8000"),
			ReadTimeout:  getdur("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getdur("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getdur("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		DB: DBConfig{
			DSN:      getenv("DB_DSN", ""),
			MaxConns: int32(getint("DB_MAX_CONNS", 10)),
			MinConns: int32(getint("DB_MIN_CONNS", 2)),
		},
		JWT: JWTConfig{
			Secret:     []byte(getenv("JWT_SECRET", "")),
			AccessTTL:  getdur("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: getdur("JWT_REFRESH_TTL", 720*time.Hour),
			Issuer:     getenv("JWT_ISSUER", "hero-typing"),
		},
		Argon2: Argon2Config{
			MemoryKiB:   uint32(getint("ARGON2_MEMORY_KIB", 65536)),
			Iterations:  uint32(getint("ARGON2_ITERATIONS", 3)),
			Parallelism: uint8(getint("ARGON2_PARALLELISM", 2)),
			SaltLen:     uint32(getint("ARGON2_SALT_LEN", 16)),
			KeyLen:      uint32(getint("ARGON2_KEY_LEN", 32)),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		},
		StaticDir: getenv("STATIC_DIR", "frontend/dist"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.DB.DSN == "" {
		return fmt.Errorf("config: DB_DSN is required")
	}
	if len(c.JWT.Secret) < 16 {
		return fmt.Errorf("config: JWT_SECRET must be at least 16 bytes")
	}
	if c.JWT.AccessTTL <= 0 || c.JWT.RefreshTTL <= 0 {
		return fmt.Errorf("config: JWT TTLs must be positive")
	}
	if c.Argon2.Iterations == 0 || c.Argon2.KeyLen == 0 || c.Argon2.SaltLen == 0 {
		return fmt.Errorf("config: invalid argon2 params")
	}
	return nil
}

func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return def
}

func getint(k string, def int) int {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getdur(k string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
