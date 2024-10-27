package providers

import (
	"context"
	"testing"
)

func TestGarantextE2E(t *testing.T) {
	p := NewGarantexProvider()
	_, err := p.GetRates(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}
