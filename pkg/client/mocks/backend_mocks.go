package mocks

import (
	"backendify/pkg/models"
	"math/rand"
	"time"
)

var mockData []models.Company = []models.Company{
	{
		ID:     "1",
		Name:   "Company A",
		Active: true,
	},
	{
		ID:     "2",
		Name:   "Company B",
		Active: true,
	},
	{
		ID:          "3",
		Name:        "Company C",
		Active:      false,
		ActiveUntil: "2023-07-15T00:00:00Z",
	},
	{
		ID:          "4",
		Name:        "Company D",
		Active:      false,
		ActiveUntil: "2023-08-01T00:00:00Z",
	},
	{
		ID:     "5",
		Name:   "Company E",
		Active: true,
	},
	{
		ID:          "6",
		Name:        "Company F",
		Active:      false,
		ActiveUntil: "2023-09-01T00:00:00Z",
	},
	{
		ID:     "7",
		Name:   "Company G",
		Active: true,
	},
	{
		ID:          "8",
		Name:        "Company H",
		Active:      false,
		ActiveUntil: "2023-07-25T00:00:00Z",
	},
	{
		ID:     "9",
		Name:   "Company I",
		Active: true,
	},
	{
		ID:     "10",
		Name:   "Company J",
		Active: true,
	},
}

type MockBackendClient struct {
}

// Returns a random Company from the mock data.
func (m MockBackendClient) FetchCompanyData(backendURL, id string) (*models.Company, error) {
	var customRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := customRand.Intn(len(mockData))

	return &mockData[randomIndex], nil
}

func (m MockBackendClient) StartWorkers() {

}
func (m MockBackendClient) StopWorkers() {

}

func (m MockBackendClient) WorkersAvailable() bool {
	return true
}
