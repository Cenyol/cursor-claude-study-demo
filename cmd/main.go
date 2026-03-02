package main

import (
	"log"
	"os"

	"user-system/internal/infrastructure/cache"
	"user-system/internal/infrastructure/persistence"
	"user-system/internal/interfaces/http"
)

func main() {
	dsn := getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/user_system?charset=utf8mb4&parseTime=True")
	redisAddr := getEnv("REDIS_ADDR", "127.0.0.1:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	db, err := persistence.NewMySQL(dsn)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	rdb, err := cache.NewRedis(redisAddr, redisPass, redisDB)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	userRepo := persistence.NewUserRepository(db)
	sessionRepo := cache.NewSessionRepository(rdb)
	handler := &http.UserHandler{UserRepo: userRepo, SessionRepo: sessionRepo}

	r := http.SetupRouter(handler)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
