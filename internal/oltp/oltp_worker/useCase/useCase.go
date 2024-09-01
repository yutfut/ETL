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

	newChan    chan models.Client
	updateChan chan models.Client

	LastInsertID uint64
	LastUpdateAT time.Time

	newTicker    *time.Ticker
	updateTicker *time.Ticker
}

func NewUseCase(
	repository repository.Repository,
	logger *log.Logger,
	dataChan chan models.Client,
	updateChan chan models.Client,
) UseCase {
	return &useCase{
		repository: repository,
		logger:     logger,

		LastInsertID: 0,
		LastUpdateAT: oltp_worker.StartTime,

		newChan:    dataChan,
		updateChan: updateChan,

		newTicker:    time.NewTicker(1 * time.Minute),
		updateTicker: time.NewTicker(1 * time.Minute),
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

			u.newChan <- item
		}

		u.LastInsertID = lastInsertID
		u.LastUpdateAT = lastUpdateAT
	}

	go u.SelectNew(ctx)

	go u.SelectUpdate(ctx)
}

func (u *useCase) SelectNew(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			u.logger.Println("SelectNew Done")
			return
		case <-u.newTicker.C:
			response, errSelectByID := u.repository.SelectByID(ctx, u.LastInsertID)
			if errSelectByID != nil {
				u.logger.Println("Start ::: u.repository.SelectByID ::: %v", errSelectByID)
			}

			var lastInsertID uint64 = 0

			for _, item := range response {
				if item.ID > lastInsertID {
					lastInsertID = item.ID
				}

				u.newChan <- item
			}

			u.LastInsertID = lastInsertID
		}
	}
}

func (u *useCase) SelectUpdate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			u.logger.Println("SelectNew Done")
			return
		case <-u.updateTicker.C:
			response, errSelectByID := u.repository.SelectByUpdateAT(ctx, u.LastUpdateAT)
			if errSelectByID != nil {
				u.logger.Println("Start ::: u.repository.SelectByID ::: %v", errSelectByID)
			}

			lastUpdateAT := oltp_worker.StartTime

			for _, item := range response {
				if item.UpdateAT.After(lastUpdateAT) {
					lastUpdateAT = item.UpdateAT
				}

				u.updateChan <- item
			}

			u.LastUpdateAT = lastUpdateAT
		}
	}
}
