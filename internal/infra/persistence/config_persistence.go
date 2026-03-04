package persistence

import (
	"context"
	"graduation-invitation/internal/domain/config"

	"github.com/jmoiron/sqlx"
)

type configPersistence struct {
	db *sqlx.DB
}

func NewConfigPersistence(db *sqlx.DB) config.ConfigRepo {
	return &configPersistence{db: db}
}

func (p *configPersistence) GetByKey(ctx context.Context, key string) (*config.Config, error) {
	var c config.Config
	query := `SELECT * FROM system_configs WHERE key = $1`
	err := p.db.GetContext(ctx, &c, query, key)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (p *configPersistence) Save(ctx context.Context, c *config.Config) error {
	query := `
		INSERT INTO system_configs (key, value, created_at, updated_at)
		VALUES (:key, :value, NOW(), NOW())
		ON CONFLICT (key) DO UPDATE
		SET value = :value, updated_at = NOW()`

	_, err := p.db.NamedExecContext(ctx, query, c)
	return err
}
