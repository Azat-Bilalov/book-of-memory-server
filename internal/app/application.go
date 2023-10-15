package app

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/dsn"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/joho/godotenv"
)

type Application struct {
	repository *repository.Repository
}

func New() (Application, error) {
	_ = godotenv.Load()
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return Application{}, err
	}

	config.Connect()

	return Application{repository: repo}, nil
}
