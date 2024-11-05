package controller

import (
	"context"
	"go-usdtrub/internal/models"
	"math/rand"
	"testing"

	pb "go-usdtrub/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpcController(t *testing.T) {
	ask, bid := rand.Float64(), rand.Float64()
	host := "localhost:8080"

	c, err := NewGrpcController(&MockService{Ask: ask, Bid: bid}, host)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	ctx := context.Background()

	// Test direct call

	rate, err := c.GetRates(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if rate.Ask != ask || rate.Bid != bid {
		t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", ask, bid, rate.Ask, rate.Bid)
	}

	// Test client call

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	client := pb.NewRatesClient(conn)
	rate, err = client.GetRates(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if rate.Ask != ask || rate.Bid != bid {
		t.Fatalf("Expected ask and bid to be %.2f and %.2f, got %.2f, %.2f", ask, bid, rate.Ask, rate.Bid)
	}
}

//// Mock Service

type MockService struct {
	Ask, Bid float64
}

func (m MockService) GetRates(ctx context.Context) (models.CurrencyRate, error) {
	return models.CurrencyRate{Ask: m.Ask, Bid: m.Bid}, nil
}
