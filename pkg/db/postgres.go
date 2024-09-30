package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	var errs error
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			errs = err
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Database Connected")
		return db, nil
	}
	panic(errs)
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate()

	if err != nil {
		panic(err)
	}
}