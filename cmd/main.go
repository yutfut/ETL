package main

import (
	"context"
	"etl/config"
	"etl/internal/oltp/generator"
	"etl/pkg/postgres"
	"log"
	"sync"
	"os/signal"
	"os"
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
				Host: s.Host,
				Port: s.Port,
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

	<-ctx.Done()
	stop()

	wg.Wait()
}