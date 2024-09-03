package usecase

import (
	"context"
	"log"
	"time"

	"etl/internal/models"
	"etl/internal/olap_worker/repository"
)

type UseCase interface{}

type useCase struct {
	clickHouse repository.ClickHouse
	logger     *log.Logger

	insertChan chan models.Client
	updateChan chan models.Client

	batch     []models.Client
	batchSize int

	insertTicker *time.Ticker
}

func NewUseCase(
	clickHouse repository.ClickHouse,
	logger *log.Logger,
	insertChan chan models.Client,
	updateChan chan models.Client,
	insertTicker time.Duration,
) UseCase {
	return &useCase{
		clickHouse:   clickHouse,
		logger:       logger,
		insertChan:   insertChan,
		updateChan:   updateChan,
		batch:        make([]models.Client, 1000),
		batchSize:    0,
		insertTicker: time.NewTicker(insertTicker),
	}
}

func (u *useCase) Insert(
	ctx context.Context,
) {
	for {
		select {
		case <-u.insertTicker.C:
			if err := u.clickHouse.Insert(ctx, u.batch); err != nil {
				u.logger.Println(err)
			} // а не сломается ли это если отменят контекст?
			u.batchSize = 0
		case item := <-u.insertChan:
			u.batch[u.batchSize] = item
			u.batchSize++

			if u.batchSize >= cap(u.batch) {
				if err := u.clickHouse.Insert(ctx, u.batch); err != nil {
					u.logger.Println(err)
				} // а не сломается ли это если отменят контекст?
				u.batchSize = 0
			}
		case <-ctx.Done():
			if err := u.clickHouse.Insert(ctx, u.batch); err != nil {
				u.logger.Println(err)
			} // а не сломается ли это если отменят контекст?

			batch := make([]models.Client, 0, len(u.insertChan))

			for item := range u.insertChan {
				batch = append(batch, item)
			}

			if err := u.clickHouse.Insert(ctx, batch); err != nil {
				u.logger.Println(err)
			} // а не сломается ли это если отменят контекст?

			u.logger.Println("insert gracefully shutdown done")
		}
	}
}

func (u *useCase) Update(
	ctx context.Context,
) {
	for {
		select {
		case item := <-u.updateChan:
			if err := u.clickHouse.Update(ctx, item); err != nil {
				u.logger.Println(err)
			} // а не сломается ли это если отменят контекст?
		case <-ctx.Done():
			for item := range u.updateChan {
				if err := u.clickHouse.Update(ctx, item); err != nil {
					u.logger.Println(err)
				} // а не сломается ли это если отменят контекст?
			}

			u.logger.Println("update gracefully shutdown done")
		}
	}
}
