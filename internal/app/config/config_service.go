package config

import (
	"context"
	"encoding/json"
	"graduation-invitation/internal/domain/config"
)

type ConfigService interface {
	GetInvitationConfig(ctx context.Context) (json.RawMessage, error)
	UpdateInvitationConfig(ctx context.Context, data json.RawMessage) error
}

type configService struct {
	repo config.ConfigRepo
}

func NewConfigService(repo config.ConfigRepo) ConfigService {
	return &configService{repo: repo}
}

const InvitationConfigKey = "invitation_theme"

func (s *configService) GetInvitationConfig(ctx context.Context) (json.RawMessage, error) {
	cfg, err := s.repo.GetByKey(ctx, InvitationConfigKey)
	if err != nil {
		// If not found, return nil (handler will handle default)
		return nil, nil
	}
	return cfg.Value, nil
}

func (s *configService) UpdateInvitationConfig(ctx context.Context, data json.RawMessage) error {
	cfg := &config.Config{
		Key:   InvitationConfigKey,
		Value: data,
	}
	return s.repo.Save(ctx, cfg)
}
