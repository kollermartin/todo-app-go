package types

import "time"

type Config struct {
	Env    string `mapstructure:"ENV"`
	Port   string `mapstructure:"PORT"`
	DBType string `mapstructure:"DB_TYPE"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort int    `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`
}

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoInput struct {
	Title string `json:"title" binding:"required"`
}
