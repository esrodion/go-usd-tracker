package controller

import (
	"log"
	"net/http"
)

type HttpController struct{}

func NewHttpController() *HttpController {
	return &HttpController{}
}

func (c *HttpController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Received health check request from", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
