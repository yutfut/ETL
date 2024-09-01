package repository

import (
	"context"
	"etl/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse interface {
	PushButch(ctx context.Context, data []models.Client) error
}

type clickHouse struct {
	driver driver.Conn
}

func NewClickHouse(
	driver driver.Conn,
) ClickHouse {
	return &clickHouse{
		driver: driver,
	}
}

func (ch *clickHouse) PushButch(
	ctx context.Context,
	request []models.Client,
) error {
	batch, err := ch.driver.PrepareBatch(
		ctx,
		insertBatch,
	)
	if err != nil {
		return err
	}

	if err = batch.AppendStruct(request); err != nil {
		return err
	}

	if err = batch.Send(); err != nil {
		return err
	}

	return nil
}

func (ch *clickHouse) Update(
	ctx context.Context,
	data models.Client,
) error {
	if err := ch.driver.Exec(
		ctx,
		update,
		data.ID,
		data.Name,
		data.Settlement,
		data.MarginAlgorithm,
		data.Gateway,
		data.Vendor,
		data.IsActive,
		data.IsPro,
		data.IsInterbank,
		data.CreateAT,
		data.UpdateAT,
	); err != nil {
		return err
	}

	return nil
}
