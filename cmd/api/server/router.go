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
	v1.HandleFunc("/product", s.baseHandler.RunAction(s.productHandler.GetProduct)).Methods("GET")
	v1.HandleFunc("/product/{id:[1-9][0-9]*}", s.baseHandler.RunAction(s.productHandler.UpdateProduct)).Methods("PUT")
	v1.HandleFunc("/product/{id:[1-9][0-9]*}", s.baseHandler.RunAction(s.productHandler.DeleteProduct)).Methods("DELETE")
	v1.HandleFunc("/product/customer", s.baseHandler.RunAction(s.productHandler.SearchSKU)).Methods("GET")
	v1.HandleFunc("/product/checkout/history", s.baseHandler.RunAction(s.checkoutHandler.GetCheckoutHistory)).Methods("GET")

}
