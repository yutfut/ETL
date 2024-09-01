package usecase

import (
	"context"
	"etl/internal/models"
	"etl/internal/oltp/oltp_worker"
	"log"
	"time"

	"etl/internal/oltp/oltp_worker/repository"
)

type UseCase interface {
}

type useCase struct {
	repository repository.Repository
	logger     *log.Logger

	dataChan chan models.Client

	LastInsertID uint64
	LastUpdateAT time.Time
}

func NewUseCase(
	repository repository.Repository,
	logger *log.Logger,
) UseCase {
	return &useCase{
		repository: repository,
		logger:     logger,

		LastInsertID: 0,
		LastUpdateAT: oltp_worker.StartTime,
	}
}

func (u *useCase) Start(ctx context.Context) {
	meta, err := u.repository.SelectMeta(ctx)
	if err != nil {
		u.logger.Println("Start ::: u.repository.SelectMeta ::: %v", err)
	}

	u.LastInsertID = meta.LastInsertID
	u.LastUpdateAT = meta.LastUpdateAT

	if u.LastInsertID == 0 {
		response, errSelectByID := u.repository.SelectByID(ctx, u.LastInsertID)
		if errSelectByID != nil {
			u.logger.Println("Start ::: u.repository.SelectByID ::: %v", errSelectByID)
		}

		var lastInsertID uint64 = 0
		lastUpdateAT := oltp_worker.StartTime

		for _, item := range response {
			if item.ID > lastInsertID {
				lastInsertID = item.ID
			}

			if item.UpdateAT.After(lastUpdateAT) {
				lastUpdateAT = item.UpdateAT
			}

			u.dataChan <- item
		}

		u.LastInsertID = lastInsertID
		u.LastUpdateAT = lastUpdateAT
	}

	go u.SelectNew(ctx)

	go u.SelectUpdate(ctx)
}

func (u *useCase) SelectNew(ctx context.Context) {}

func (u *useCase) SelectUpdate(ctx context.Context) {}
