package generator

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Generator interface {
	Generate(ctx context.Context, wg *sync.WaitGroup)
}

type generator struct {
	pgx    *pgxpool.Pool
	logger *log.Logger
}

func NewGenerator(
	pgx *pgxpool.Pool,
	logger *log.Logger,
) Generator {
	return &generator{
		pgx:    pgx,
		logger: logger,
	}
}

func (g *generator) Generate(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	updateTicker := time.NewTicker(5 * time.Second)
	i := 0

	for {
		select {
		case <-ctx.Done():
			g.logger.Println("Generator done")
			return
		default:
			select {
			case <-updateTicker.C:
				if _, err := g.pgx.Exec(
					context.Background(),
					update,
					gofakeit.Name(),
					rand.IntN(i),
				); err != nil {
					g.logger.Println(err)
				}
			default:
				if _, err := g.pgx.Exec(
					context.Background(),
					insert,
					gofakeit.Name(),
					gofakeit.Currency().Short,
					gofakeit.Uint8(),
					gofakeit.Bool(),
					gofakeit.Bool(),
					gofakeit.Bool(),
					gofakeit.Bool(),
					gofakeit.Bool(),
				); err != nil {
					g.logger.Println(err)
				} else {
					i++
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
