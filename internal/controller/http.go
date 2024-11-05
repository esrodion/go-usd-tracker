package controller

import (
	"encoding/json"
	"go-usdtrub/internal/traces"
	"go-usdtrub/pkg/logger"
	"net/http"
)

type HttpController struct {
	service Service
}

// @title           Go USDT/RUB tracker
// @version         1.0
//
// @license.name  CC0
//
// @host      localhost:8081
// @BasePath  /
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewHttpController(service Service) *HttpController {
	return &HttpController{service: service}
}

// HealthCheck godoc
// @Summary      Check service's health
// @Tags         General
// @Success      200  {string}  "ok"
// @Router       /healthy [get]
func (c *HttpController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := traces.Start(r.Context(), "HttpHealthCheck")
	defer span.End()

	logger.Logger().Sugar().Named("HTTP").Debug("Received health check request from ", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// GetRates godoc
// @Summary      USDT/RUB pair current rate
// @Tags         Rates
// @Produce      json
// @Success      200  {object}  models.CurrencyRate
// @Failure      500  {string}  ""
// @Router       /rates [get]
func (c *HttpController) GetRates(w http.ResponseWriter, r *http.Request) {
	ctx, span := traces.Start(r.Context(), "HttpGetRates")
	defer span.End()

	rate, err := c.service.GetRates(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
