package usecase

import (
	"context"
	"etl/internal/models"
	"etl/internal/oltp/oltp_worker"
	"log"
	"sync"
	"time"

	"etl/internal/oltp/oltp_worker/repository"
)

type UseCase interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type useCase struct {
	postgreSQLID string

	repository repository.Repository
	logger     *log.Logger

	insertChan chan models.OLAPClient
	updateChan chan models.OLAPClient

	LastInsertID uint64
	LastUpdateAT time.Time

	newTicker    *time.Ticker
	updateTicker *time.Ticker
}

func NewUseCase(
	postgreSQLID string,
	repository repository.Repository,
	logger *log.Logger,
	insertChan chan models.OLAPClient,
	updateChan chan models.OLAPClient,
) UseCase {
	return &useCase{
		postgreSQLID: postgreSQLID,

		repository: repository,
		logger:     logger,

		LastInsertID: 0,
		LastUpdateAT: oltp_worker.StartTime,

		insertChan: insertChan,
		updateChan: updateChan,

		newTicker:    time.NewTicker(1 * time.Minute),
		updateTicker: time.NewTicker(1 * time.Minute),
	}
}

func (u *useCase) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	meta, err := u.repository.SelectMeta(
		context.Background(),
	)
	if err != nil {
		u.logger.Println("Start ::: u.repository.SelectMeta ::: %v", err)
	}

	u.LastInsertID = meta.LastInsertID
	u.LastUpdateAT = meta.LastUpdateAT

	if u.LastInsertID == 0 {
		response, errSelectByID := u.repository.SelectByID(
			context.Background(),
			u.LastInsertID,
		)
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

			u.insertChan <- models.OLAPClient{
				PostgreSQLID:    u.postgreSQLID,
				ID:              item.ID,
				Name:            item.Name,
				Settlement:      item.Settlement,
				MarginAlgorithm: item.MarginAlgorithm,
				Gateway:         item.Gateway,
				Vendor:          item.Vendor,
				IsActive:        item.IsActive,
				IsPro:           item.IsPro,
				IsInterbank:     item.IsInterbank,
				CreateAT:        item.CreateAT,
				UpdateAT:        item.UpdateAT,
			}
		}

		u.LastInsertID = lastInsertID
		u.LastUpdateAT = lastUpdateAT

		_, errUpdateMeta := u.repository.UpdateMeta(
			context.Background(),
			repository.ReadMeta{
				LastInsertID: lastInsertID,
				LastUpdateAT: lastUpdateAT,
			},
		)
		if errUpdateMeta != nil {
			u.logger.Println("Start ::: u.repository.UpdateMeta ::: %v", errUpdateMeta)
		}

		//u.LastInsertID = updateMeta.LastInsertID
		//u.LastUpdateAT = updateMeta.LastUpdateAT
	}

	wg.Add(2)

	go u.insert(ctx, wg)

	go u.update(ctx, wg)
}

func (u *useCase) insert(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			u.logger.Println("insert Done")
			return
		case <-u.newTicker.C:
			response, errSelectByID := u.repository.SelectByID(
				context.Background(),
				u.LastInsertID,
			)
			if errSelectByID != nil {
				u.logger.Println("Start ::: u.repository.SelectByID ::: %v", errSelectByID)
			}

			var lastInsertID uint64 = 0

			for _, item := range response {
				if lastInsertID < item.ID {
					lastInsertID = item.ID
				}

				u.insertChan <- models.OLAPClient{
					PostgreSQLID:    u.postgreSQLID,
					ID:              item.ID,
					Name:            item.Name,
					Settlement:      item.Settlement,
					MarginAlgorithm: item.MarginAlgorithm,
					Gateway:         item.Gateway,
					Vendor:          item.Vendor,
					IsActive:        item.IsActive,
					IsPro:           item.IsPro,
					IsInterbank:     item.IsInterbank,
					CreateAT:        item.CreateAT,
					UpdateAT:        item.UpdateAT,
				}
			}

			u.LastInsertID = lastInsertID

			_, errUpdateMetaLastInsertID := u.repository.UpdateMetaLastInsertID(
				context.Background(),
				repository.ReadMeta{
					LastInsertID: lastInsertID,
				},
			)
			if errUpdateMetaLastInsertID != nil {
				u.logger.Println("Start ::: u.repository.UpdateMetaLastInsertID ::: %v", errUpdateMetaLastInsertID)
			}

			//u.LastInsertID = UpdateMetaLastInsertID.LastInsertID
		}
	}
}

func (u *useCase) update(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			u.logger.Println("update Done")
			return
		case <-u.updateTicker.C:
			response, errSelectByID := u.repository.SelectByUpdateAT(
				context.Background(),
				u.LastUpdateAT,
			)
			if errSelectByID != nil {
				u.logger.Println("update ::: u.repository.SelectByUpdateAT ::: %v", errSelectByID)
			}

			if len(response) == 0 {
				continue
			}

			lastUpdateAT := oltp_worker.StartTime

			for _, item := range response {
				if item.UpdateAT.After(lastUpdateAT) {
					lastUpdateAT = item.UpdateAT
				}

				u.updateChan <- models.OLAPClient{
					PostgreSQLID:    u.postgreSQLID,
					ID:              item.ID,
					Name:            item.Name,
					Settlement:      item.Settlement,
					MarginAlgorithm: item.MarginAlgorithm,
					Gateway:         item.Gateway,
					Vendor:          item.Vendor,
					IsActive:        item.IsActive,
					IsPro:           item.IsPro,
					IsInterbank:     item.IsInterbank,
					CreateAT:        item.CreateAT,
					UpdateAT:        item.UpdateAT,
				}
			}

			u.LastUpdateAT = lastUpdateAT

			_, err := u.repository.UpdateMetaLastUpdateAT(
				context.Background(),
				repository.ReadMeta{
					LastUpdateAT: lastUpdateAT,
				},
			)
			if err != nil {
				u.logger.Println("update ::: u.repository.UpdateMetaLastUpdateAT ::: %v", err)
				return
			}

			//u.LastUpdateAT = UpdateMetaLastUpdateAT.LastUpdateAT
		}
	}
}
