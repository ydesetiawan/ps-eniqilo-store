package server

import (
	"github.com/alexliesenfeld/health"
	healthcheckhandler "ps-cats-social/internal/healthcheck/handler"
)

func (s *Server) setupRouter() {
	v1 := s.router.PathPrefix("/v1").Subrouter().StrictSlash(true)
	v1.HandleFunc("/health", health.NewHandler(healthcheckhandler.HealthCheck())).Methods("GET")

}
