package config

import (
	"context"
	"encoding/json"
	"time"
)

type Config struct {
	Key       string          `db:"key" json:"key"`
	Value     json.RawMessage `db:"value" json:"value"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type ConfigRepo interface {
	GetByKey(ctx context.Context, key string) (*Config, error)
	Save(ctx context.Context, config *Config) error
}
