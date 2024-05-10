package server

import (
	healthcheckhandler "ps-eniqilo-store/internal/healthcheck/handler"

	"github.com/alexliesenfeld/health"
)

func (s *Server) setupRouter() {
	v1 := s.router.PathPrefix("/v1").Subrouter().StrictSlash(true)
	v1.HandleFunc("/health", health.NewHandler(healthcheckhandler.HealthCheck())).Methods("GET")

	v1.HandleFunc("/staff/register", s.baseHandler.RunAction(s.userHandler.Register)).Methods("POST")
	v1.HandleFunc("/staff/login", s.baseHandler.RunAction(s.userHandler.Login)).Methods("POST")

	v1.HandleFunc("/customer/register", s.baseHandler.RunActionAuth(s.customerHandler.CreateCustomer)).Methods("POST")
	v1.HandleFunc("/customer", s.baseHandler.RunActionAuth(s.customerHandler.GetCustomers)).Methods("GET")

	v1.HandleFunc("/product", s.baseHandler.RunActionAuth(s.productHandler.CreateProduct)).Methods("POST")
	v1.HandleFunc("/product", s.baseHandler.RunActionAuth(s.productHandler.GetProduct)).Methods("GET")
	v1.HandleFunc("/product/{id:[1-9][0-9]*}", s.baseHandler.RunActionAuth(s.productHandler.UpdateProduct)).Methods("PUT")
	v1.HandleFunc("/product/{id:[1-9][0-9]*}", s.baseHandler.RunActionAuth(s.productHandler.DeleteProduct)).Methods("DELETE")
	v1.HandleFunc("/product/customer", s.baseHandler.RunActionAuth(s.productHandler.SearchSKU)).Methods("GET")
	v1.HandleFunc("/product/checkout/history", s.baseHandler.RunActionAuth(s.checkoutHandler.GetCheckoutHistory)).Methods("GET")

}
