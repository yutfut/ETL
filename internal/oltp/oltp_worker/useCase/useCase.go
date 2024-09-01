package usecase

import (
	"context"
	"log"
	"time"

	"etl/internal/oltp/oltp_worker/repository"
)

type UseCase interface {

}

type useCase struct {
	repository repository.Repository
	logger *log.Logger

	lastInsertID uint64
	lastUpdateAT time.Time
}

func NewUseCase(
	repository repository.Repository,
	logger *log.Logger,
) UseCase {
	return &useCase{
		repository: repository,
		logger: logger,

		LastInsertID: 0,
		LastUpdateAT: time.Date(1970, 01, 01, 0, 0, 0),
	}
}

func (u *useCase) Start(ctx context.Context) {
	_, err := u.repository.SelectMeta(ctx)
	if err != nil {
		u.logger.Println("Start ::: u.repository.SelectMeta ::: %v", err)
	}
}

