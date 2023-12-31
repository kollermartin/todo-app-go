package types

import (
	"time"
)

type Config struct {
	Env            string `mapstructure:"ENV"`
	Port           string `mapstructure:"PORT"`
	DBType         string `mapstructure:"DB_TYPE"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPass         string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`
}

type Todo struct {
	ID         int       `json:"id"`
	ExternalID string    `json:"external_id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
}

type TodoResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoInput struct {
	Title string `json:"title" binding:"required"`
}
