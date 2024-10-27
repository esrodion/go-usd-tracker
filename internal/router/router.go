package router

import (
	"go-usdtrub/internal/controller"
	"net/http"
)

func NewRouter(c *controller.HttpController) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("GET /healthy", c.HealthCheck)

	return m
}
