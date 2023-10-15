package config

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/dsn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// todo: перенести в конфиг
var (
	DB  *gorm.DB
	err error
)

func Connect() {
	DB, err = gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Database connected!")
}
