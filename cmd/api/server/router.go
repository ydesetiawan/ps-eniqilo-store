package server

import (
	"github.com/alexliesenfeld/health"
	healthcheckhandler "ps-eniqilo-store/internal/healthcheck/handler"
)

func (s *Server) setupRouter() {
	v1 := s.router.PathPrefix("/v1").Subrouter().StrictSlash(true)
	v1.HandleFunc("/health", health.NewHandler(healthcheckhandler.HealthCheck())).Methods("GET")

	//TODO change to auth
	v1.HandleFunc("/product", s.baseHandler.RunAction(s.productHandler.CreateProduct)).Methods("POST")

}
