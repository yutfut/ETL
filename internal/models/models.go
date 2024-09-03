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
	PostgreSQLID    string    `db:"postgresql_id"`
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
