package repository

import (
	"chess-backend/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type IRepositories interface {
	GetAllSteps() ([]domain.InfoStep, error)
}

type Repositories struct {
	IRepositories
}

func NewRepositories(mongoDB *mongo.Client) *Repositories {
	return &Repositories{
		NewRepoStep(mongoDB),
	}
}
