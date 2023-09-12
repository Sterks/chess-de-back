package service

import (
	"chess-backend/internal/domain"
	"chess-backend/internal/repository"
)

type IProcessingService interface {
	ReadProcessing(string) error
}

type IServicesSteps interface {
	GetSteps() ([]domain.InfoStep, error)
}

type Services struct {
	StepsService      IServicesSteps
	ProcessingService IProcessingService
}

type Deps struct {
	Repos repository.IRepositories
}

func NewServices(deps Deps) *Services {
	return &Services{
		ProcessingService: NewFileProcessingService(deps.Repos),
		StepsService:      NewStepsService(deps.Repos),
	}
}
