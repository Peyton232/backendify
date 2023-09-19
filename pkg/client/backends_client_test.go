package client

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestIsClosedDateInThePast(t *testing.T) {
	t.Run("Date in the past", func(t *testing.T) {
		pastDate := time.Now().Add(-time.Hour * 24)
		dateStr := pastDate.Format(time.RFC3339)
		result := isClosedDateInThePast(dateStr)
		assert.True(t, result, "Expected date in the past to be true")
	})

	t.Run("Date in the future", func(t *testing.T) {
		futureDate := time.Now().Add(time.Hour * 24)
		dateStr := futureDate.Format(time.RFC3339)
		result := isClosedDateInThePast(dateStr)
		assert.False(t, result, "Expected date in the future to be false")
	})

	t.Run("Invalid date", func(t *testing.T) {
		dateStr := "invalid_date"
		result := isClosedDateInThePast(dateStr)
		assert.False(t, result, "Expected invalid date to be false")
	})

	t.Run("Empty date", func(t *testing.T) {
		dateStr := ""
		result := isClosedDateInThePast(dateStr)
		assert.False(t, result, "Expected empty date to be false")
	})
}

func TestParseCompanyResponse(t *testing.T) {
	t.Run("Unsupported Content Type", func(t *testing.T) {
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/unsupported")

		_, err := ParseCompanyResponse(resp, "123")
		assert.Error(t, err, "Expected an error for unsupported content type")
	})

	t.Run("V1: No ClosedOn Date", func(t *testing.T) {
		v1Response := struct {
			CN        string `json:"cn"`
			CreatedOn string `json:"created_on"`
			ClosedOn  string `json:"closed_on,omitempty"`
		}{
			CN:        "Company Name",
			CreatedOn: "2023-01-01T00:00:00Z", // Date in the past
		}
		respBody, _ := json.Marshal(v1Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v1")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.True(t, company.Active, "Expected company to be active")
		assert.Empty(t, company.ActiveUntil, "Expected ActiveUntil to be empty")
	})

	t.Run("V1: ClosedOn Date in the Past", func(t *testing.T) {
		v1Response := struct {
			CN        string `json:"cn"`
			CreatedOn string `json:"created_on"`
			ClosedOn  string `json:"closed_on,omitempty"`
		}{
			CN:        "Company Name",
			CreatedOn: "2023-01-01T00:00:00Z", // Date in the past
			ClosedOn:  "2023-02-01T00:00:00Z", // Date in the past
		}
		respBody, _ := json.Marshal(v1Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v1")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.False(t, company.Active, "Expected company to be inactive")
		assert.Equal(t, v1Response.ClosedOn, company.ActiveUntil, "Expected ActiveUntil to be set")
	})

	t.Run("V1: ClosedOn Date in the Future", func(t *testing.T) {
		v1Response := struct {
			CN        string `json:"cn"`
			CreatedOn string `json:"created_on"`
			ClosedOn  string `json:"closed_on,omitempty"`
		}{
			CN:        "Company Name",
			CreatedOn: "2023-01-01T00:00:00Z", // Date in the past
			ClosedOn:  "2124-01-01T00:00:00Z", // Date in the future
		}
		respBody, _ := json.Marshal(v1Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v1")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.True(t, company.Active, "Expected company to be active")
		assert.Equal(t, v1Response.ClosedOn, company.ActiveUntil, "Expected ActiveUntil to be set")
	})

	t.Run("V2: No DissolvedOn Date", func(t *testing.T) {
		v2Response := struct {
			CompanyName string `json:"company_name"`
			TIN         string `json:"tin"`
			DissolvedOn string `json:"dissolved_on,omitempty"`
		}{
			CompanyName: "Company Name",
			TIN:         "TIN123",
		}
		respBody, _ := json.Marshal(v2Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v2")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.True(t, company.Active, "Expected company to be active")
		assert.Empty(t, company.ActiveUntil, "Expected ActiveUntil to be empty")
	})

	t.Run("V2: DissolvedOn Date in the Past", func(t *testing.T) {
		v2Response := struct {
			CompanyName string `json:"company_name"`
			TIN         string `json:"tin"`
			DissolvedOn string `json:"dissolved_on,omitempty"`
		}{
			CompanyName: "Company Name",
			TIN:         "TIN123",
			DissolvedOn: "2023-01-01T00:00:00Z", // Date in the past
		}
		respBody, _ := json.Marshal(v2Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v2")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.False(t, company.Active, "Expected company to be inactive")
		assert.Equal(t, v2Response.DissolvedOn, company.ActiveUntil, "Expected ActiveUntil to be set")
	})

	t.Run("V2: DissolvedOn Date in the Future", func(t *testing.T) {
		v2Response := struct {
			CompanyName string `json:"company_name"`
			TIN         string `json:"tin"`
			DissolvedOn string `json:"dissolved_on,omitempty"`
		}{
			CompanyName: "Company Name",
			TIN:         "TIN123",
			DissolvedOn: "2124-01-01T00:00:00Z", // Date in the future
		}
		respBody, _ := json.Marshal(v2Response)
		resp := &fasthttp.Response{}
		resp.Header.Set("Content-Type", "application/x-company-v2")
		resp.SetBody(respBody)

		company, err := ParseCompanyResponse(resp, "123")
		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, company, "Expected company data")
		assert.True(t, company.Active, "Expected company to be active")
		assert.Equal(t, v2Response.DissolvedOn, company.ActiveUntil, "Expected ActiveUntil to be set")
	})
}
