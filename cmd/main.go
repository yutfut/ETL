package main

import (
	"context"
	"etl/config"
	"etl/internal/models"
	olapRepository "etl/internal/olap_worker/repository"
	olapUseCase "etl/internal/olap_worker/useCase"
	"etl/internal/oltp/generator"
	oltpRepository "etl/internal/oltp/oltp_worker/repository"
	oltpUseCase "etl/internal/oltp/oltp_worker/useCase"
	clickhouse "etl/pkg/click_house"
	"etl/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	logger := log.New(
		os.Stdout,
		"ERR\t",
		log.Ldate|log.Ltime|log.Llongfile,
	)

	secret, err := config.LoadConfig("./config/config.json")
	if err != nil {
		logger.Fatalln(err)
	}

	wg := &sync.WaitGroup{}

	for _, s := range secret.Postgres {
		postgresPool, err := postgres.NewPostgresPool(
			ctx,
			postgres.Postgres{
				Host:     s.Host,
				Port:     s.Port,
				Username: s.Username,
				Password: s.Password,
				Database: s.Database,
			},
		)
		if err != nil {
			logger.Fatal(err)
		}

		wg.Add(1)
		go generator.NewGenerator(
			postgresPool,
			logger,
		).Generate(ctx, wg)
	}

	insertChan := make(chan models.OLAPClient, 10000)
	updateChan := make(chan models.OLAPClient, 10000)

	for _, s := range secret.Postgres {
		postgresPool, err := postgres.NewPostgresPool(
			ctx,
			postgres.Postgres{
				Host:     s.Host,
				Port:     s.Port,
				Username: s.Username,
				Password: s.Password,
				Database: s.Database,
			},
		)
		if err != nil {
			logger.Fatal(err)
		}

		postgresRepository := oltpRepository.NewRepository(
			postgresPool,
			logger,
		)

		ptUseCase := oltpUseCase.NewUseCase(
			s.PostgreSQLID,
			postgresRepository,
			logger,
			insertChan,
			updateChan,
		)

		ptUseCase.Start(ctx, wg)
	}

	clickhouseConn, err := clickhouse.Connect(
		clickhouse.ClickHouse{
			Host:     secret.ClickHouse.Host,
			Port:     secret.ClickHouse.Port,
			Database: secret.ClickHouse.Database,
			Username: secret.ClickHouse.Username,
			Password: secret.ClickHouse.Password,
			Debug:    secret.ClickHouse.Debug,
		},
	)
	if err != nil {
		logger.Fatal(err)
	}

	clickhouseRepository := olapRepository.NewClickHouse(
		clickhouseConn,
		logger,
	)

	apUseCase := olapUseCase.NewUseCase(
		clickhouseRepository,
		logger,
		insertChan,
		updateChan,
	)

	wg.Add(2)

	go apUseCase.Insert(ctx, wg)

	go apUseCase.Update(ctx, wg)

	<-ctx.Done()
	stop()

	wg.Wait()

	log.Println("main done")
}
