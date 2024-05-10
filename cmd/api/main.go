package main

import (
	"fmt"
	stdlog "log"
	"os"
	"ps-eniqilo-store/cmd/api/server"
	"ps-eniqilo-store/configs"

	customerhandler "ps-eniqilo-store/internal/customer/handler"
	customerrepository "ps-eniqilo-store/internal/customer/repository"
	customerservice "ps-eniqilo-store/internal/customer/service"

	userhandler "ps-eniqilo-store/internal/user/handler"
	userrepository "ps-eniqilo-store/internal/user/repository"
	userservice "ps-eniqilo-store/internal/user/service"

	checkouthandler "ps-eniqilo-store/internal/checkout/handler"
	checkoutrepository "ps-eniqilo-store/internal/checkout/repository"
	checkoutservice "ps-eniqilo-store/internal/checkout/service"

	producthandler "ps-eniqilo-store/internal/product/handler"
	productrepository "ps-eniqilo-store/internal/product/repository"
	productservice "ps-eniqilo-store/internal/product/service"

	"ps-eniqilo-store/internal/shared"
	bhandler "ps-eniqilo-store/pkg/base/handler"
	"ps-eniqilo-store/pkg/logger"
	psqlqgen "ps-eniqilo-store/pkg/psqlqgen"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var port int

var httpCmd = &cobra.Command{
	Use:   "http [OPTIONS]",
	Short: "Run HTTP API",
	Long:  "Run HTTP API for SCM",
	RunE:  runHttpCommand,
}

var (
	params          map[string]string
	baseHandler     *bhandler.BaseHTTPHandler
	userHandler     *userhandler.UserHandler
	customerHandler *customerhandler.CustomerHandler
	productHandler  *producthandler.ProductHandler
	checkoutHandler *checkouthandler.CheckoutHandler
)

func init() {
	httpCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the HTTP server")
}

func main() {
	if err := httpCmd.Execute(); err != nil {
		slog.Error(fmt.Sprintf("Error on command execution: %s", err.Error()))
		os.Exit(1)
	}
}

func logLevel() slog.Level {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initLogger() {
	{

		log, err := logger.SlogOption{
			Resource: map[string]string{
				"service.name":        shared.ServiceName,
				"service.ns":          "eniqilo-store",
				"service.instance_id": "random-uuid",
				"service.version":     "v.0",
				"service.env":         "staging",
			},
			ContextExtractor:   nil,
			AttributeFormatter: nil,
			Writer:             os.Stdout,
			Leveler:            logLevel(),
		}.NewSlog()
		if err != nil {
			err = fmt.Errorf("prepare logger error: %w", err)
			stdlog.Fatal(err) // if logger cannot be prepared (commonly due to option value error), use std logger.
			return
		}

		// Set logger as global logger.
		slog.SetDefault(log)
	}
}

func runHttpCommand(cmd *cobra.Command, args []string) error {
	initLogger()
	initInfra()

	httpServer := server.NewServer(
		baseHandler,
		userHandler,
		customerHandler,
		productHandler,
		checkoutHandler,
		port,
	)

	return httpServer.Run()
}

func dbInitConnection() *sqlx.DB {
	return psqlqgen.Init(configs.Init(), shared.ServiceName)
}

func initInfra() {
	db := dbInitConnection()

	customerRepository := customerrepository.NewCustomerRepositoryImpl(db)
	customerService := customerservice.NewCustomerServiceImpl(customerRepository)
	customerHandler = customerhandler.NewCustomerHandler(customerService)

	userRepository := userrepository.NewUserRepositoryImpl(db)
	userService := userservice.NewUserServiceImpl(userRepository)
	userHandler = userhandler.NewUserHandler(userService)

	productRepository := productrepository.NewProductRepositoryImpl(db)
	productService := productservice.NewProductServiceImpl(productRepository)
	productHandler = producthandler.NewProductHandler(productService)

	checkoutRepository := checkoutrepository.NewCheckoutRepositoryImpl(db)
	checkoutService := checkoutservice.NewCheckoutServiceImpl(checkoutRepository)
	checkoutHandler = checkouthandler.NewCheckoutHandler(checkoutService)
}
