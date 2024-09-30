package app

import (
	pg "ChatApp/pkg/db"

	"gorm.io/gorm"
)

func (app *App) Run() {
	var db *gorm.DB
	var err error
	db, err = pg.NewDB(app.config.DBHost)

	if err != nil {
		panic(err)
	}

	pg.Migrate(db)
}

func initControllers() {
}