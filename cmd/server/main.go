package main

import (
	"ChatApp/configs"
	App "ChatApp/internal/app"
	"fmt"
)

func main() {
	config, err := configs.LoadConfig()

	if err != nil {
		panic(err)
	}
	fmt.Println(config.AppName)

	app := App.NewApp(config)
	app.Run()
}