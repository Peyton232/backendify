package api

import (
	"backendify/pkg/models"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func (cr *CustomRouter) Status(ctx *fasthttp.RequestCtx) {
	// Check if your solution is ready to accept requests
	if !cr.IsReadyToAcceptRequests() {
		// If not ready, return a 503 status code (Service Unavailable)
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}

	// If ready, return a 200/OK status code
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/plain")
	ctx.Write([]byte("OK"))
}

// Define a function to check if your solution is ready
func (cr *CustomRouter) IsReadyToAcceptRequests() bool {
	return cr.BackendClient.WorkersAvailable()
}

func (cr *CustomRouter) GetCompany(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))
	iso := string(ctx.QueryArgs().Peek("country_iso"))

	// Check if either 'id' or 'country_iso' is missing
	if id == "" || iso == "" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// Check if ISO code is associated with a backend
	backend, found := cr.Backends[iso]
	if !found {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// Use a buffered channel to communicate the response
	ch := make(chan struct {
		company *models.Company
		err     error
	}, 1)

	// Acquire a semaphore before starting the goroutine
	if err := cr.fetchSemaphore.Acquire(ctx, 1); err != nil {
		ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
		return
	}
	defer cr.fetchSemaphore.Release(1)

	// Start a goroutine to fetch company data concurrently
	go func() {
		defer close(ch) // Close the channel when done
		company, err := cr.BackendClient.FetchCompanyData(backend, id)
		ch <- struct {
			company *models.Company
			err     error
		}{company, err}
	}()

	// Wait for the goroutine to finish and send the response
	result := <-ch
	if result.err != nil {
		cr.Logger.Error("An error occurred:", result.err)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	if result.company == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// Respond with company data
	cr.Logger.Info("Company data retrieved successfully")
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	json.NewEncoder(ctx).Encode(result.company)
}
