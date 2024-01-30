package config

import "github.com/spf13/viper"

type App struct {
	Name string
	Env  string
	Port string
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Type     string
}

type Config struct {
	App *App
	Db  *Db
}

func New() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	app := &App{
		Name: viper.GetString("APP_NAME"),
		Env:  viper.GetString("ENV"),
		Port: viper.GetString("PORT"),
	}

	db:= &Db{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASS"),
		Name:     viper.GetString("DB_NAME"),
		Type:     viper.GetString("DB_TYPE"),
	}

	return &Config{
		App: app,
		Db:  db,
	}, nil
}
