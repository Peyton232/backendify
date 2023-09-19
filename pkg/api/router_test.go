package api

import (
	"backendify/pkg/config"
	"backendify/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewRouter(t *testing.T) {
	// Create a sample BackendConfig
	backends := config.BackendConfig{
		"us": "http://localhost:9001",
		"ru": "http://localhost:9002",
	}

	config := models.Config{
		Application: models.ApplicationConfig{
			MockFlag: true,
		},
		Limiter: models.LimiterConfig{
			Limit:  100,
			Period: "15s",
		},
	}

	// Create a new CustomRouter with mockFlag set to true
	router, err := NewRouter(backends, &config, nil)
	assert.Nil(t, err)

	// Create a sample fasthttp request for testing
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/status")
	ctx.Request.Header.SetMethod("GET")

	// Create a fasthttp response for capturing the response
	resp := &fasthttp.Response{}

	// Use the router to handle the request
	router.HandleRequest(ctx)

	// Check if the response status code is as expected (200 OK)
	assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}
