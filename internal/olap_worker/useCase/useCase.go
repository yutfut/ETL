package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"etl/internal/models"
	"etl/internal/olap_worker/repository"
)

type UseCase interface {
	Insert(ctx context.Context, wg *sync.WaitGroup)
	Update(ctx context.Context, wg *sync.WaitGroup)
}

type useCase struct {
	clickHouse repository.ClickHouse
	logger     *log.Logger

	insertChan chan models.OLAPClient
	updateChan chan models.OLAPClient

	batch     []models.OLAPClient
	batchSize int

	insertTicker *time.Ticker
}

func NewUseCase(
	clickHouse repository.ClickHouse,
	logger *log.Logger,
	insertChan chan models.OLAPClient,
	updateChan chan models.OLAPClient,
) UseCase {
	return &useCase{
		clickHouse:   clickHouse,
		logger:       logger,
		insertChan:   insertChan,
		updateChan:   updateChan,
		batch:        make([]models.OLAPClient, 1000),
		batchSize:    0,
		insertTicker: time.NewTicker(1 * time.Minute),
	}
}

func (u *useCase) Insert(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			if err := u.clickHouse.Insert(
				context.Background(),
				u.batch[0:u.batchSize],
			); err != nil {
				u.logger.Println(err)
			}

			batch := make([]models.OLAPClient, 0, len(u.insertChan))

			for item := range u.insertChan {
				batch = append(batch, item)
			}

			if err := u.clickHouse.Insert(
				context.Background(),
				u.batch[0:u.batchSize],
			); err != nil {
				u.logger.Println(err)
			}

			u.logger.Println("insert gracefully shutdown done")
		default:
			select {
			case <-u.insertTicker.C:
				if err := u.clickHouse.Insert(
					context.Background(), //todo: thinking
					u.batch[0:u.batchSize],
				); err != nil {
					u.logger.Println(err)
				}
				u.batchSize = 0
			case item := <-u.insertChan:
				u.batch[u.batchSize] = item
				u.batchSize++

				if u.batchSize >= cap(u.batch) {
					if err := u.clickHouse.Insert(
						context.Background(), //todo: thinking
						u.batch[0:u.batchSize],
					); err != nil {
						u.logger.Println(err)
					}
					u.batchSize = 0
				}
			}
		}
	}
}

func (u *useCase) Update(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			for item := range u.updateChan {
				if err := u.clickHouse.Update(
					context.Background(),
					item,
				); err != nil {
					u.logger.Println(err)
				}
			}

			u.logger.Println("update gracefully shutdown done")
		default:
			select {
			case item := <-u.updateChan:
				if err := u.clickHouse.Update(
					context.Background(),
					item,
				); err != nil {
					u.logger.Println(err)
					u.updateChan <- item
				}
			}
		}
	}
}
