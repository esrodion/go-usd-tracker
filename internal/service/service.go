package service

import (
	"context"
	"fmt"
	"go-usdtrub/internal/models"
	"go-usdtrub/internal/service/providers"
	"go-usdtrub/internal/traces"
)

type Repository interface {
	AddRates(context.Context, models.CurrencyRate) error
}

type RatesProvider interface {
	GetRates(context.Context) (models.CurrencyRate, error)
}

type Service struct {
	repo     Repository
	provider RatesProvider
}

type option func(s *Service)

func WithProvider(provider RatesProvider) option {
	return func(s *Service) {
		s.provider = provider
	}
}

func NewService(repo Repository, opts ...option) *Service {
	s := &Service{repo: repo}
	for _, opt := range opts {
		opt(s)
	}

	if s.provider == nil {
		s.provider = providers.NewGarantexProvider()
	}

	return s
}

func (s *Service) GetRates(ctx context.Context) (models.CurrencyRate, error) {
	ctx, span := traces.Start(ctx, "ServiceGetRates")
	defer span.End()

	rate, err := s.provider.GetRates(ctx)
	if err != nil {
		return rate, fmt.Errorf("service.Service.GetRates: %w", err)
	}

	err = s.repo.AddRates(ctx, rate)
	if err != nil {
		return rate, fmt.Errorf("service.Service.GetRates: %w", err)
	}

	return rate, nil
}
