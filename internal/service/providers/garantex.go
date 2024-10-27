package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-usdtrub/internal/models"
	"io"
	"net/http"
	"time"
)

const GarantexEndpoint string = "https://garantex.org/api/v2/depth?market=usdtrub"

type GarantexProvider struct{}

func NewGarantexProvider() *GarantexProvider {
	return &GarantexProvider{}
}

func (p *GarantexProvider) GetRates(ctx context.Context) (models.CurrenceyRate, error) {
	result := models.CurrenceyRate{}

	req, err := http.NewRequestWithContext(ctx, "GET", GarantexEndpoint, nil)
	if err != nil {
		return result, fmt.Errorf("service.Service.getGarantexRates: could not create http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("service.Service.getGarantexRates: could not process http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("service.Service.getGarantexRates: bad response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("service.Service.getGarantexRates: could not read response body: %w", err)
	}

	depth, err := unmarshalDepth(data)
	if err != nil {
		return result, fmt.Errorf("service.Service.getGarantexRates: could not parse response body to json: %w", err)
	}

	result.Timestamp = time.UnixMicro(depth.Timestamp)

	if len(depth.Asks) > 0 {
		result.Ask = depth.Asks[0].Price
	}

	if len(depth.Bids) > 0 {
		result.Bid = depth.Bids[0].Price
	}

	return result, nil
}

//// Service

func unmarshalDepth(data []byte) (depth, error) {
	var r depth
	err := json.Unmarshal(data, &r)
	return r, err
}

type depth struct {
	Timestamp int64        `json:"timestamp"`
	Asks      []depthEntry `json:"asks"`
	Bids      []depthEntry `json:"bids"`
}

type depthEntry struct {
	Price  float64 `json:"price,string"`
	Volume float64 `json:"volume,string"`
	Amount float64 `json:"amount,string"`
	Factor float64 `json:"factor,string"`
	Type   string  `json:"type"`
}
