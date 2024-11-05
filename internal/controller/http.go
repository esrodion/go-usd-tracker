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

func NewHttpController(service Service) *HttpController {
	return &HttpController{service: service}
}

func (c *HttpController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := traces.Start(r.Context(), "HttpHealthCheck")
	defer span.End()

	logger.Logger().Sugar().Named("HTTP").Debug("Received health check request from ", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

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
