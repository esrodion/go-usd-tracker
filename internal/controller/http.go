package controller

import (
	"go-usdtrub/pkg/logger"
	"net/http"
)

type HttpController struct{}

func NewHttpController() *HttpController {
	return &HttpController{}
}

func (c *HttpController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Logger().Sugar().Named("HTTP").Debug("Received health check request from ", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
