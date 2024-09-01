package repository

import (
	"context"
	"etl/internal/models"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SelectMeta(ctx context.Context) (ReadMeta, error)
	UpdateMeta(ctx context.Context, request ReadMeta) (ReadMeta, error)
	SelectByID(ctx context.Context, request uint64) ([]models.Client, error)
	SelectByUpdateAT(ctx context.Context, request time.Time) ([]models.Client, error)
}

type repository struct {
	pgx    *pgxpool.Pool
	logger *log.Logger
}

func NewRepository(
	pgx *pgxpool.Pool,
	logger *log.Logger,
) Repository {
	return &repository{
		pgx:    pgx,
		logger: logger,
	}
}

type ReadMeta struct {
	LastInsertID uint64    `db:"last_insert_id"`
	LastUpdateAT time.Time `db:"last_update_at"`
}

func (r *repository) SelectMeta(
	ctx context.Context,
) (ReadMeta, error) {
	rows, err := r.pgx.Query(
		ctx,
		selectMeta,
	)
	if err != nil {
		r.logger.Printf("ReadMeta ::: r.pgx.Query ::: %v", err)
		return ReadMeta{}, err
	}

	response, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[ReadMeta],
	)
	if err != nil {
		r.logger.Printf("ReadMeta ::: pgx.CollectOneRow ::: %v", err)
		return ReadMeta{}, err
	}

	return response, nil
}

func (r *repository) UpdateMeta(
	ctx context.Context,
	request ReadMeta,
) (ReadMeta, error) {
	rows, err := r.pgx.Query(
		ctx,
		updateMeta,
		request.LastInsertID,
		request.LastUpdateAT,
	)
	if err != nil {
		r.logger.Printf("UpdateMeta ::: r.pgx.Query ::: %v", err)
		return ReadMeta{}, err
	}

	response, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[ReadMeta],
	)
	if err != nil {
		r.logger.Printf("UpdateMeta ::: pgx.CollectOneRow ::: %v", err)
		return ReadMeta{}, err
	}

	return response, nil
}

func (r *repository) SelectByID(
	ctx context.Context,
	request uint64,
) (
	[]models.Client, error,
) {
	rows, err := r.pgx.Query(
		ctx,
		selectByID,
		request,
	)
	if err != nil {
		r.logger.Printf("SelectInsert ::: r.pgx.Query ::: %v", err)
		return []models.Client{}, err
	}

	response, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[models.Client],
	)
	if err != nil {
		r.logger.Printf("UpdateMeta ::: pgx.CollectOneRow ::: %v", err)
		return []models.Client{}, err
	}

	return response, nil
}

func (r *repository) SelectByUpdateAT(
	ctx context.Context,
	request time.Time,
) (
	[]models.Client, error,
) {
	rows, err := r.pgx.Query(
		ctx,
		selectByUpdateAT,
		request,
	)
	if err != nil {
		r.logger.Printf("SelectInsert ::: r.pgx.Query ::: %v", err)
		return []models.Client{}, err
	}

	response, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[models.Client],
	)
	if err != nil {
		r.logger.Printf("UpdateMeta ::: pgx.CollectOneRow ::: %v", err)
		return []models.Client{}, err
	}

	return response, nil
}
