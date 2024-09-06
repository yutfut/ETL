package main

import (
	"context"
	"etl/config"
	"etl/internal/models"
	"etl/internal/olap_worker/repository"
	clickhouse "etl/pkg/click_house"
	"github.com/brianvoe/gofakeit"
	"log"
	"os"
)

func main() {
	logger := log.New(
		os.Stdout,
		"ERR\t",
		log.Ldate|log.Ltime|log.Llongfile,
	)

	secret, err := config.LoadConfig("./config/config.json")
	if err != nil {
		logger.Fatalln(err)
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

	batch, err := clickhouseConn.PrepareBatch(
		context.Background(),
		repository.InsertBatch,
	)
	if err != nil {
		logger.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		m := create()

		if err = batch.AppendStruct(
			&m,
		); err != nil {
			logger.Println(err)
			return
		}
	}

	if err = batch.Send(); err != nil {
		logger.Println(err)
		return
	}

	logger.Println("batch sent")
}

func create() models.OLAPClient {
	return models.OLAPClient{
		PostgreSQLID:    gofakeit.Name(),
		ID:              gofakeit.Uint64(),
		Name:            gofakeit.Name(),
		Settlement:      gofakeit.Currency().Short,
		MarginAlgorithm: gofakeit.Uint8(),
		Gateway:         gofakeit.Bool(),
		Vendor:          gofakeit.Bool(),
		IsActive:        gofakeit.Bool(),
		IsPro:           gofakeit.Bool(),
		IsInterbank:     gofakeit.Bool(),
		CreateAT:        gofakeit.Date(),
		UpdateAT:        gofakeit.Date(),
	}
}
