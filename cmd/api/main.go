package main

import (
	"fmt"
	stdlog "log"
	"os"
	"ps-eniqilo-store/cmd/api/server"
	producthandler "ps-eniqilo-store/internal/product/handler"
	productrepository "ps-eniqilo-store/internal/product/repository"
	productservice "ps-eniqilo-store/internal/product/service"
	"ps-eniqilo-store/internal/shared"
	bhandler "ps-eniqilo-store/pkg/base/handler"
	"ps-eniqilo-store/pkg/logger"
	psqlqgen "ps-eniqilo-store/pkg/psqlqgen"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
	params         map[string]string
	baseHandler    *bhandler.BaseHTTPHandler
	productHandler *producthandler.ProductHandler
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
		baseHandler, productHandler, port,
	)

	return httpServer.Run()
}

func dbInitConnection() *sqlx.DB {
	godotenv.Load(".env")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	uname := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbparams := os.Getenv("DB_PARAMS")

	return psqlqgen.Init(host, port, uname, pass, dbname, dbparams, shared.ServiceName)
}

func initInfra() {
	db := dbInitConnection()

	productRepository := productrepository.NewProductRepositoryImpl(db)
	productService := productservice.NewProductServiceImpl(productRepository)
	productHandler = producthandler.NewProductHandler(productService)

}
