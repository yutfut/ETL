package repository

import (
	"context"
	"etl/internal/models"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse interface {
	Insert(ctx context.Context, data []models.OLAPClient) error
	Update(ctx context.Context, request models.OLAPClient) error
}

type clickHouse struct {
	driver driver.Conn
	logger *log.Logger
}

func NewClickHouse(
	driver driver.Conn,
	logger *log.Logger,
) ClickHouse {
	return &clickHouse{
		driver: driver,
		logger: logger,
	}
}

func (ch *clickHouse) Insert(
	ctx context.Context,
	request []models.OLAPClient,
) error {
	batch, err := ch.driver.PrepareBatch(
		ctx,
		InsertBatch,
	)
	if err != nil {
		ch.logger.Println(err)
		return err
	}

	for _, item := range request {
		if err = batch.AppendStruct(
			&item,
		); err != nil {
			ch.logger.Println(err)
			return err
		}
	}

	if err = batch.Send(); err != nil {
		ch.logger.Println(err)
		return err
	}

	return nil
}

func (ch *clickHouse) Update(
	ctx context.Context,
	request models.OLAPClient,
) error {
	if err := ch.driver.Exec(
		ctx,
		update,
		request.PostgreSQLID,
		request.ID,
		request.Name,
		request.Settlement,
		request.MarginAlgorithm,
		request.Gateway,
		request.Vendor,
		request.IsActive,
		request.IsPro,
		request.IsInterbank,
		request.UpdateAT,
	); err != nil {
		ch.logger.Println(err)
		return err
	}

	return nil
}
