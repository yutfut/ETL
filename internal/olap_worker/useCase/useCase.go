package usecase

import (
	"context"
	"etl/internal/models"
	"etl/internal/olap_worker/repository"
	"log"
	"time"
)

type UseCase interface{}

type useCase struct {
	clickHouse repository.ClickHouse
	dataChan   chan models.Client
	batch      []models.Client
	batchSize   int
	pushTicker *time.Ticker
}

func NewUseCase(
	clickHouse repository.ClickHouse,
	dataChan chan models.Client,
	tickerTime time.Duration,
) UseCase {
	return &useCase{
		clickHouse: clickHouse,
		dataChan:   dataChan,
		batch:      make([]models.Client, 1000),
		batchSize:   0,
		pushTicker: time.NewTicker(tickerTime),
	}
}

func (u *useCase) Start(ctx context.Context) {
	for {
		select {
		case <-u.pushTicker.C:
			if err := u.clickHouse.PushButch(ctx, u.batch); err != nil {
				log.Println(err)
			} // а не сломается ли это если отменят контекст?
			u.batchSize = 0
		case item := <- u.dataChan:
			u.batch[u.batchSize] = item
			u.batchSize++

			if u.batchSize >= cap(u.batch) {
				if err := u.clickHouse.PushButch(ctx, u.batch); err != nil {
					log.Println(err)
				} // а не сломается ли это если отменят контекст?
				u.batchSize = 0
			}
		case <- ctx.Done():
			if err := u.clickHouse.PushButch(ctx, u.batch); err != nil {
				log.Println(err)
			} // а не сломается ли это если отменят контекст?

			batch := make([]models.Client, 0, len(u.dataChan))

			for item := range u.dataChan {
				batch = append(batch, item)
			}

			if err := u.clickHouse.PushButch(ctx, batch); err != nil {
				log.Println(err)
			} // а не сломается ли это если отменят контекст?
		}
	}
}
