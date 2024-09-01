package usecase

import (
	"context"

	"etl/internal/oltp/oltp_worker/repository"
)

type UseCase interface {

}

type useCase struct {
	repository repository.Repository
	logger *log.Logger
}

func NewUseCase(
	repository repository.Repository,
	logger *log.Logger,
) UseCase {
	return &useCase{
		repository: repository,
		logger: logger,
	}
}

func (u *useCase) Start(ctx context.Context) {
	_, err := u.repository.SelectMeta(ctx)
	if err != nil {
		logger.Println("Start ::: u.repository.SelectMeta ::: %v", err)
	}
}

