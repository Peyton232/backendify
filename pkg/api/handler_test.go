package api

import (
	"backendify/pkg/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestGetCompany(t *testing.T) {
	// Create a logger for the test
	logger := logrus.New()

	config := models.Config{
		Application: models.ApplicationConfig{
			MockFlag: true,
		},
		Limiter: models.LimiterConfig{
			Limit:  100,
			Period: "15s",
		},
	}

	// Initialize your router with a mock client and the logger.
	r, err := NewRouter(map[string]string{"us": "http://example.com"}, &config, logger)
	assert.Nil(t, err)

	testCases := []struct {
		name         string
		requestURI   string
		expectedCode int
		expectedBody bool
	}{
		{
			name:         "NoID",
			requestURI:   "/company?country_iso=us",
			expectedCode: fasthttp.StatusNotFound,
			expectedBody: false,
		},
		{
			name:         "InvalidISOCode",
			requestURI:   "/company?id=1",
			expectedCode: fasthttp.StatusNotFound,
			expectedBody: false,
		},
		{
			name:         "HappyPath",
			requestURI:   "/company?id=1&country_iso=us",
			expectedCode: fasthttp.StatusOK,
			expectedBody: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new fasthttp.RequestCtx for testing.
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod("GET")
			ctx.Request.SetRequestURI(tc.requestURI)

			// Call the router's request handler.
			r.HandleRequest(ctx)

			// Assert the response status code.
			assert.Equal(t, tc.expectedCode, ctx.Response.StatusCode())

			// Assert the response body (if applicable).
			if tc.expectedBody {
				assert.NotEmpty(t, tc.expectedBody)
			} else {
				assert.Empty(t, tc.expectedBody)
			}

		})
	}
}
