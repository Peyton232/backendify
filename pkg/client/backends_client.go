package client

import (
	"backendify/pkg/models"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/valyala/fasthttp"
)

// CompanyFetcher is an interface for fetching company data.
type CompanyFetcher interface {
	FetchCompanyData(backendURL, id string) (*models.Company, error)
	StartWorkers()
	StopWorkers()
	WorkersAvailable() bool
}

type BackendClient struct {
	httpClient       *fasthttp.Client
	requestPool      sync.Pool
	cache            *lru.Cache
	requests         chan requestInfo
	workers          int
	availableWorkers int
	wg               sync.WaitGroup
}

type requestInfo struct {
	backendURL string
	id         string
	result     chan<- *models.Company
}

var (
	ErrCacheMiss       = errors.New("cache miss")
	ErrInvalidResponse = errors.New("invalid response")
)

// NewBackendClient initializes a new BackendClient with the given cache size and worker count.
func NewBackendClient(cacheSize, workerCount int) (*BackendClient, error) {
	httpClient := &fasthttp.Client{}
	requestPool := sync.Pool{
		New: func() interface{} {
			return new(fasthttp.Request)
		},
	}

	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	return &BackendClient{
		httpClient:  httpClient,
		requestPool: requestPool,
		cache:       cache,
		requests:    make(chan requestInfo),
		workers:     workerCount,
	}, nil
}

// StartWorkers starts the worker goroutines to process requests.
func (bc *BackendClient) StartWorkers() {
	for i := 0; i < bc.workers; i++ {
		bc.wg.Add(1)
		go bc.worker()
		bc.availableWorkers++
	}
}

// StopWorkers stops the worker goroutines.
func (bc *BackendClient) StopWorkers() {
	close(bc.requests)
	bc.wg.Wait()
}

func (bc *BackendClient) WorkersAvailable() bool {
	return bc.availableWorkers > 0
}

// FetchCompanyData sends a request to the worker pool to fetch company data.
func (bc *BackendClient) FetchCompanyData(backendURL, id string) (*models.Company, error) {
	resultChan := make(chan *models.Company)
	req := requestInfo{
		backendURL: backendURL,
		id:         id,
		result:     resultChan,
	}

	select {
	case bc.requests <- req:
		result := <-resultChan
		close(req.result)
		return result, nil
	default:
		// Handle case where sending the request to the worker pool fails (e.g., worker pool is full)
		close(req.result)
		return nil, errors.New("failed to send request to worker pool")
	}
}

func (bc *BackendClient) worker() {
	defer bc.wg.Done()
	for {
		req, more := <-bc.requests
		bc.availableWorkers--
		if !more {
			// The requests channel is closed, so the worker can exit.
			return
		}
		cachedResponse, found := bc.cache.Get(req.id)
		var company *models.Company

		// Extract backendURL from the request
		backendURL := req.backendURL
		id := req.id
		if found {
			if cachedCompany, ok := cachedResponse.(*models.Company); ok {
				company = cachedCompany
			}
		} else {
			request := bc.requestPool.Get().(*fasthttp.Request)
			request.SetRequestURI(backendURL + "/companies/" + id)

			var resp fasthttp.Response
			if err := bc.httpClient.Do(request, &resp); err != nil {
				request.Header.Reset()
				bc.requestPool.Put(request)

				req.result <- nil
				bc.availableWorkers--
				continue
			}

			parsedCompany, err := ParseCompanyResponse(&resp, id)
			if err != nil {
				request.Header.Reset()
				bc.requestPool.Put(request)

				req.result <- nil
				bc.availableWorkers--
				continue
			}
			company = parsedCompany
			bc.cache.Add(id, company)
			request.Header.Reset()
			bc.requestPool.Put(request)
		}

		req.result <- company
		bc.availableWorkers--

	}
}

func isClosedDateInThePast(dateStr string) bool {
	if dateStr == "" {
		return false // No date provided, assume active
	}

	// Parse the date string to a time.Time
	closedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return false // Error parsing date, assume active
	}

	// Compare the parsed date with the current time
	return closedDate.Before(time.Now())
}

// ParseCompanyResponse parses the response based on content type and returns a Company.
func ParseCompanyResponse(resp *fasthttp.Response, id string) (*models.Company, error) {
	var company models.Company

	contentType := string(resp.Header.ContentType())

	switch {
	case strings.HasPrefix(contentType, "application/x-company-v1"):
		var v1Response struct {
			CN        string `json:"cn"`
			CreatedOn string `json:"created_on"`
			ClosedOn  string `json:"closed_on,omitempty"`
		}
		if err := json.Unmarshal(resp.Body(), &v1Response); err != nil {
			return nil, err
		}
		company.ID = id
		company.Name = v1Response.CN
		company.Active = v1Response.ClosedOn == "" || !isClosedDateInThePast(v1Response.ClosedOn)
		company.ActiveUntil = v1Response.ClosedOn

	case strings.HasPrefix(contentType, "application/x-company-v2"):
		var v2Response struct {
			CompanyName string `json:"company_name"`
			TIN         string `json:"tin"`
			DissolvedOn string `json:"dissolved_on,omitempty"`
		}
		if err := json.Unmarshal(resp.Body(), &v2Response); err != nil {
			return nil, err
		}
		company.ID = id
		company.Name = v2Response.CompanyName
		company.Active = v2Response.DissolvedOn == "" || !isClosedDateInThePast(v2Response.DissolvedOn)
		company.ActiveUntil = v2Response.DissolvedOn

	default:
		return nil, errors.New("unsupported content type")
	}

	return &company, nil
}
