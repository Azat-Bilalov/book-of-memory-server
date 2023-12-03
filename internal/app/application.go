package app

import (
	"context"
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/dsn"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/redis"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/joho/godotenv"
)

type Application struct {
	config     *config.Config
	repository *repository.Repository
	redis      *redis.Client
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	_ = godotenv.Load()
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}
	log.Println("app repo", repo)

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Application{
		config:     cfg,
		repository: repo,
		redis:      redisClient,
	}, nil
}
