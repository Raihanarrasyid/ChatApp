package app

import "ChatApp/configs"

type App struct {
	config *configs.Config
}

func NewApp(config *configs.Config) *App {
	return &App{config}
}
