package router

import (
	"fmt"
	"go-usdtrub/internal/config"
	"go-usdtrub/internal/controller"
	"go-usdtrub/internal/docs"
	"go-usdtrub/internal/metrics"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(c *controller.HttpController, cfg *config.Config) *http.ServeMux {

	m := http.NewServeMux()
	m.HandleFunc("GET /healthy", metrics.WrapHandlerFunc("health_check", "/healthy", c.HealthCheck))
	m.HandleFunc("GET /rates", metrics.WrapHandlerFunc("get_rates", "/rates", c.GetRates))
	m.Handle("GET /metrics", metrics.HandlerHTTP())

	docs.SwaggerInfo.Host = "localhost:" + cfg.HttpPort
	m.Handle("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", docs.SwaggerInfo.Host)),
	))

	return m
}
