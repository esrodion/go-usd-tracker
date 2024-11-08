package service

import (
	"context"
	"go-usdtrub/internal/models"
	"testing"
)

func TestService(t *testing.T) {
	s := NewService(&MockRepository{}, WithProvider(&MockRatesProvider{}))
	_, err := s.GetRates(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}

//// Mocks

type MockRepository struct{}

func (repo *MockRepository) AddRates(ctx context.Context, rate models.CurrencyRate) error {
	return nil
}

type MockRatesProvider struct{}

func (rp *MockRatesProvider) GetRates(ctx context.Context) (models.CurrencyRate, error) {
	return models.CurrencyRate{}, nil
}
