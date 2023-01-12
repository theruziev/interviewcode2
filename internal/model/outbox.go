package model

import (
	"time"
)

type OutboxStatus string

const (
	CreatedStatus = OutboxStatus("created")
	DoneStatus    = OutboxStatus("done")
	ErrorStatus   = OutboxStatus("error")
)

type OutBox struct {
	ID        int64        `db:"id"`
	Topic     string       `db:"topic"`
	Data      any          `db:"msg"`
	Status    OutboxStatus `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
