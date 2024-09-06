package models

import "time"

type OLTPClient struct {
	ID              uint64    `db:"id"`
	Name            string    `db:"name"`
	Settlement      string    `db:"settlement"`
	MarginAlgorithm uint8     `db:"margin_algorithm"`
	Gateway         bool      `db:"gateway"`
	Vendor          bool      `db:"vendor"`
	IsActive        bool      `db:"is_active"`
	IsPro           bool      `db:"is_pro"`
	IsInterbank     bool      `db:"is_interbank"`
	CreateAT        time.Time `db:"create_at"`
	UpdateAT        time.Time `db:"update_at"`
}

type OLAPClient struct {
	PostgreSQLID    string    `ch:"postgresql_id"`
	ID              uint64    `ch:"id"`
	Name            string    `ch:"name"`
	Settlement      string    `ch:"settlement"`
	MarginAlgorithm uint8     `ch:"margin_algorithm"`
	Gateway         bool      `ch:"gateway"`
	Vendor          bool      `ch:"vendor"`
	IsActive        bool      `ch:"is_active"`
	IsPro           bool      `ch:"is_pro"`
	IsInterbank     bool      `ch:"is_interbank"`
	CreateAT        time.Time `ch:"create_at"`
	UpdateAT        time.Time `ch:"update_at"`
}
