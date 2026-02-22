// Package khaleejiapi provides the official Go SDK for KhaleejiAPI,
// the MENA region's developer API platform.
//
// Usage:
//
//	client := khaleejiapi.New("kapi_live_your_key")
//	result, err := client.Validation.ValidateEmail(ctx, "user@example.com")
package khaleejiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://khaleejiapi.dev/api/v1"
	defaultTimeout = 30 * time.Second
	userAgent      = "khaleejiapi-go/1.0.0"
	version        = "1.0.0"
)

// Config holds the configuration for the KhaleejiAPI client.
type Config struct {
	// APIKey is your API key from https://khaleejiapi.dev/dashboard/api-keys
	APIKey string
	// BaseURL is the API base URL (default: https://khaleejiapi.dev/api/v1)
	BaseURL string
	// Timeout is the HTTP request timeout (default: 30s)
	Timeout time.Duration
	// MaxRetries is the maximum retry attempts on rate limiting (default: 2)
	MaxRetries int
	// HTTPClient allows providing a custom http.Client
	HTTPClient *http.Client
}

// Client is the KhaleejiAPI SDK client.
type Client struct {
	config     Config
	httpClient *http.Client

	// Resource namespaces
	Validation    *ValidationResource
	Geo           *GeoResource
	Finance       *FinanceResource
	Communication *CommunicationResource
	Islamic       *IslamicResource
	Utility       *UtilityResource
}

// New creates a new KhaleejiAPI client with the given API key.
func New(apiKey string) *Client {
	return NewWithConfig(Config{APIKey: apiKey})
}

// NewWithConfig creates a new KhaleejiAPI client with custom configuration.
func NewWithConfig(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 2
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.Timeout}
	}

	c := &Client{
		config:     config,
		httpClient: httpClient,
	}

	c.Validation = &ValidationResource{client: c}
	c.Geo = &GeoResource{client: c}
	c.Finance = &FinanceResource{client: c}
	c.Communication = &CommunicationResource{client: c}
	c.Islamic = &IslamicResource{client: c}
	c.Utility = &UtilityResource{client: c}

	return c
}

// apiResponse is the standard API response wrapper.
type apiResponse[T any] struct {
	Data T             `json:"data"`
	Meta *ResponseMeta `json:"meta,omitempty"`
}

// ResponseMeta contains metadata from API responses.
type ResponseMeta struct {
	Timestamp string `json:"timestamp,omitempty"`
	Cached    bool   `json:"cached,omitempty"`
}

// RateLimitInfo contains rate limit header information.
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     int64
}

// APIError represents an error returned by the KhaleejiAPI.
type APIError struct {
	StatusCode    int
	Code          string
	Message       string
	RateLimitInfo *RateLimitInfo
}

func (e *APIError) Error() string {
	return fmt.Sprintf("khaleejiapi: %s (%d): %s", e.Code, e.StatusCode, e.Message)
}

type apiErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// get performs a GET request and decodes the response into T.
func doGet[T any](c *Client, ctx context.Context, path string, params map[string]string) (T, error) {
	var zero T

	u, err := url.Parse(c.config.BaseURL + path)
	if err != nil {
		return zero, &APIError{StatusCode: 0, Code: "INVALID_URL", Message: err.Error()}
	}

	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return zero, &APIError{StatusCode: 0, Code: "REQUEST_ERROR", Message: err.Error()}
	}

	return execute[T](c, req)
}

// doPost performs a POST request and decodes the response into T.
func doPost[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	var zero T

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return zero, &APIError{StatusCode: 0, Code: "MARSHAL_ERROR", Message: err.Error()}
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+path, bodyReader)
	if err != nil {
		return zero, &APIError{StatusCode: 0, Code: "REQUEST_ERROR", Message: err.Error()}
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return execute[T](c, req)
}

// execute runs the HTTP request with retry logic.
func execute[T any](c *Client, req *http.Request) (T, error) {
	var zero T
	var lastErr error

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			time.Sleep(backoff)

			// Clone request for retry
			newReq := req.Clone(req.Context())
			if req.Body != nil {
				if body, ok := req.Body.(io.Seeker); ok {
					body.Seek(0, io.SeekStart)
				}
			}
			req = newReq
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = &APIError{StatusCode: 0, Code: "NETWORK_ERROR", Message: err.Error()}
			continue
		}
		defer resp.Body.Close()

		rateLimitInfo := parseRateLimitHeaders(resp)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var apiResp apiResponse[T]
			if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
				return zero, &APIError{StatusCode: resp.StatusCode, Code: "DECODE_ERROR", Message: err.Error()}
			}
			return apiResp.Data, nil
		}

		if resp.StatusCode == 429 {
			lastErr = &APIError{
				StatusCode:    429,
				Code:          "RATE_LIMITED",
				Message:       fmt.Sprintf("Rate limited. Retry after %d seconds", rateLimitInfo.Reset),
				RateLimitInfo: rateLimitInfo,
			}
			if attempt < c.config.MaxRetries {
				continue
			}
			return zero, lastErr
		}

		// Parse error response
		bodyBytes, _ := io.ReadAll(resp.Body)
		var apiErr apiErrorResponse
		code := "SERVER_ERROR"
		message := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if json.Unmarshal(bodyBytes, &apiErr) == nil && apiErr.Error.Code != "" {
			code = apiErr.Error.Code
			message = apiErr.Error.Message
		}

		switch resp.StatusCode {
		case 401:
			return zero, &APIError{StatusCode: 401, Code: "UNAUTHORIZED", Message: "Invalid or missing API key"}
		case 403:
			return zero, &APIError{StatusCode: 403, Code: "FORBIDDEN", Message: message}
		case 404:
			return zero, &APIError{StatusCode: 404, Code: "NOT_FOUND", Message: "Resource not found"}
		default:
			return zero, &APIError{StatusCode: resp.StatusCode, Code: code, Message: message}
		}
	}

	if lastErr != nil {
		return zero, lastErr
	}
	return zero, &APIError{StatusCode: 500, Code: "UNKNOWN", Message: "Unknown error"}
}

func parseRateLimitHeaders(resp *http.Response) *RateLimitInfo {
	info := &RateLimitInfo{}
	if v := resp.Header.Get("X-RateLimit-Limit"); v != "" {
		info.Limit, _ = strconv.Atoi(v)
	}
	if v := resp.Header.Get("X-RateLimit-Remaining"); v != "" {
		info.Remaining, _ = strconv.Atoi(v)
	}
	if v := resp.Header.Get("X-RateLimit-Reset"); v != "" {
		info.Reset, _ = strconv.ParseInt(v, 10, 64)
	}
	return info
}
