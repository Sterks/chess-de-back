package service

import (
	"chess-backend/internal/domain"
	"chess-backend/internal/repository"
)

type StepService struct {
	Repo repository.IRepositories
}

func NewStepsService(repo repository.IRepositories) *StepService {
	return &StepService{
		Repo: repo,
	}
}

func (s *StepService) GetSteps() ([]domain.InfoStep, error) {
	allSteps, err := s.Repo.GetAllSteps()
	if err != nil {
		return nil, err
	}
	return allSteps, err
}
