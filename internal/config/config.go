package config

import "os"

// Config holds simple configuration for the application
type Config struct {
	Port        string
	APIKey      string
	CouponFiles []string
}

// Load creates a new Config with environment variables or defaults
func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "8080"),
		APIKey: getEnv("API_KEY", "apitest"),
		CouponFiles: []string{
			getEnv("COUPON_FILE1", "couponbase1.txt"),
			getEnv("COUPON_FILE2", "couponbase2.txt"),
			getEnv("COUPON_FILE3", "couponbase3.txt"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
