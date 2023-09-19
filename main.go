package main

import (
	"backendify/pkg/api"
	"backendify/pkg/config"
	"backendify/pkg/models"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	// Set the maximum number of CPUs to utilize
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Load configuration and initialize logger
	backends, appConfig, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	logger := initializeLogger()

	// Create router and server
	router, err := api.NewRouter(backends, appConfig, logger)
	if err != nil {
		log.Fatal(err)
	}
	server := createServer(router.HandleRequest, appConfig)

	// Use a channel to listen for OS interrupt signals (e.g., Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	startServer(server, appConfig, logger)

	// Start a goroutine to handle the second interrupt signal
	go func() {
		<-interrupt // Wait for the first interrupt signal
		logger.Info("Received SIGINT. Press Ctrl+C again to force shutdown.")
		<-interrupt // Wait for the second interrupt signal
		logger.Warn("Received second SIGINT. Forcing program shutdown...")
		os.Exit(1) // Exit with an error code to indicate a forced shutdown
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	<-interrupt

	// Shutdown the server gracefully and stop worker pool
	shutdownServer(server, logger)
	duration := 3 * time.Second
	time.Sleep(duration)
	router.ShutDown()
}

func loadConfiguration() (config.BackendConfig, *models.Config, error) {
	args := os.Args[1:]
	backends := config.LoadBackends(args)
	appConfig, err := config.LoadConfig()
	return backends, appConfig, err
}

func initializeLogger() *logrus.Logger {
	logger := logrus.New()
	// Configure your logger settings here, if needed
	return logger
}

func createServer(router fasthttp.RequestHandler, appConfig *models.Config) *fasthttp.Server {
	server := &fasthttp.Server{
		Handler:           router,
		ReadTimeout:       appConfig.Server.ReadTimeout,
		WriteTimeout:      appConfig.Server.WriteTimeout,
		ReduceMemoryUsage: appConfig.Server.ReduceMemoryUsage,
		GetOnly:           appConfig.Server.GetOnly,
	}
	return server
}

func startServer(server *fasthttp.Server, appConfig *models.Config, logger *logrus.Logger) {
	go func() {
		logger.Infof("Starting server on %d...\n", appConfig.Server.Port)
		err := server.ListenAndServe(":" + strconv.Itoa(appConfig.Server.Port))
		if err != nil {
			logger.Errorf("Error: %s\n", err)
		}
	}()
}

func shutdownServer(server *fasthttp.Server, logger *logrus.Logger) {
	logger.Info("Shutting down server gracefully...")
	err := server.Shutdown()
	if err != nil {
		logger.Error("Error shutting down server: ", err)
	} else {
		logger.Info("Server has been shut down.")
	}
}
