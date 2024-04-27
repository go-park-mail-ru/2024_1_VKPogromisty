package repository

import (
	customtime "socio/pkg/time"
)

type CSAT struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewCSAT(db DBPool, tp customtime.TimeProvider) *CSAT {
	return &CSAT{
		db: db,
		TP: tp,
	}
}
