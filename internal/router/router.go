package router

import (
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/metrics"
	"net/http"
)

func NewRouter(c *controller.HttpController) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("GET /healthy", metrics.WrapHandlerFunc("health_check", "/healthy", c.HealthCheck))
	m.Handle("GET /metrics", metrics.HandlerHTTP())

	return m
}
