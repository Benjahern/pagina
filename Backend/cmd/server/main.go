package main

import (
	"fmt"
	"log"

	"tienda/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	fmt.Println("✓ Config loaded successfully")
	fmt.Println()

	fmt.Printf("Server:    port=%s, gin_mode=%s, log_level=%s, seed=%v\n",
		cfg.Server.Port, cfg.Server.GinMode, cfg.Server.LogLevel, cfg.Server.Seed)
	fmt.Printf("           CORS allowed: %v\n", cfg.Server.CORSAllowed)

	fmt.Printf("Database:  host=%s, port=%d, user=%s, name=%s\n",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name)
	fmt.Printf("           password set: %v (len=%d)\n",
		cfg.Database.Password != "", len(cfg.Database.Password))

	fmt.Printf("JWT:       secret_len=%d bytes (min 32)\n", len(cfg.JWT.Secret))
	fmt.Printf("           expiration=%dh, bcrypt_cost=%d\n",
		cfg.JWT.ExpirationHours, cfg.JWT.BcryptCost)
}
