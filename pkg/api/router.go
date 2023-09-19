package api

import (
	"backendify/pkg/client"
	"backendify/pkg/client/mocks"
	"backendify/pkg/config"
	"backendify/pkg/models"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/semaphore"
)

// CustomRouter extends the fasthttp.RequestHandler
type CustomRouter struct {
	Backends       config.BackendConfig
	BackendClient  client.CompanyFetcher
	Logger         *logrus.Logger
	fetchSemaphore *semaphore.Weighted
}

func NewRouter(backends config.BackendConfig, config *models.Config, logger *logrus.Logger) (*CustomRouter, error) {
	r := &CustomRouter{
		Backends:       backends,
		Logger:         logger,
		fetchSemaphore: semaphore.NewWeighted(100),
	}

	// if mock mode is on setup mock client
	if config.Application.MockFlag {
		r.BackendClient = mocks.MockBackendClient{}
		return r, nil
	}

	// Create a new instance of BackendClient with worker pool support
	newClient, err := client.NewBackendClient(config.Application.CacheSize, config.Application.Workers)
	if err != nil {
		return nil, err
	}
	newClient.StartWorkers()
	r.BackendClient = newClient

	return r, nil
}

func (cr *CustomRouter) HandleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/status":
		LoggingMiddleware(cr.Status)(ctx)
	case "/company":
		LoggingMiddleware(cr.GetCompany)(ctx)
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

func (cr *CustomRouter) ShutDown() {
	cr.BackendClient.StopWorkers()
}
