package powerbiapi

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RetryConfig represents configuration for retry behavior
type RetryConfig struct {
	MaxRetries      int           // Maximum number of retry attempts
	InitialDelay    time.Duration // Initial delay before first retry
	MaxDelay        time.Duration // Maximum delay between retries
	BackoffFactor   float64       // Exponential backoff factor
	JitterFactor    float64       // Random jitter factor (0-1)
	RetryableStatus []int         // HTTP status codes to retry on
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      5,
		InitialDelay:    1 * time.Second,
		MaxDelay:        60 * time.Second,
		BackoffFactor:   2.0,
		JitterFactor:    0.3,
		RetryableStatus: []int{429, 502, 503, 504}, // Rate limit and server errors
	}
}

// EnhancedRetryRoundTripper implements enhanced retry logic with exponential backoff
type EnhancedRetryRoundTripper struct {
	innerRoundTripper http.RoundTripper
	config            *RetryConfig
}

// NewEnhancedRetryRoundTripper creates a new enhanced retry round tripper
func NewEnhancedRetryRoundTripper(next http.RoundTripper, config *RetryConfig) http.RoundTripper {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &EnhancedRetryRoundTripper{
		innerRoundTripper: next,
		config:            config,
	}
}

// RoundTrip implements the http.RoundTripper interface with enhanced retry logic
func (rt *EnhancedRetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error
	
	for attempt := 0; attempt <= rt.config.MaxRetries; attempt++ {
		// Clone the request for each attempt
		reqCopy := rt.cloneRequest(req)
		
		resp, err := rt.innerRoundTripper.RoundTrip(reqCopy)
		lastResp = resp
		lastErr = err
		
		// If successful or non-retryable error, return immediately
		if err != nil || !rt.shouldRetry(resp, attempt) {
			return resp, err
		}
		
		// Calculate delay for next retry
		delay := rt.calculateDelay(resp, attempt)
		
		// Wait before retrying
		select {
		case <-time.After(delay):
			// Continue to next retry
		case <-req.Context().Done():
			// Request context cancelled
			return nil, req.Context().Err()
		}
	}
	
	return lastResp, lastErr
}

// shouldRetry determines if a response should trigger a retry
func (rt *EnhancedRetryRoundTripper) shouldRetry(resp *http.Response, attempt int) bool {
	if attempt >= rt.config.MaxRetries {
		return false
	}
	
	if resp == nil {
		return false
	}
	
	// Check if status code is retryable
	for _, status := range rt.config.RetryableStatus {
		if resp.StatusCode == status {
			return true
		}
	}
	
	return false
}

// calculateDelay calculates the delay before the next retry attempt
func (rt *EnhancedRetryRoundTripper) calculateDelay(resp *http.Response, attempt int) time.Duration {
	var delay time.Duration
	
	// For 429 responses, check Retry-After header first
	if resp != nil && resp.StatusCode == 429 {
		if retryAfter := rt.readRetryAfterHeader(resp); retryAfter > 0 {
			// Respect the Retry-After header
			delay = retryAfter
		}
	}
	
	// If no Retry-After header, use exponential backoff
	if delay == 0 {
		delay = rt.calculateExponentialBackoff(attempt)
	}
	
	// Apply jitter to prevent thundering herd
	delay = rt.applyJitter(delay)
	
	// Ensure delay doesn't exceed maximum
	if delay > rt.config.MaxDelay {
		delay = rt.config.MaxDelay
	}
	
	return delay
}

// calculateExponentialBackoff calculates exponential backoff delay
func (rt *EnhancedRetryRoundTripper) calculateExponentialBackoff(attempt int) time.Duration {
	backoff := float64(rt.config.InitialDelay) * math.Pow(rt.config.BackoffFactor, float64(attempt))
	return time.Duration(backoff)
}

// applyJitter adds random jitter to the delay
func (rt *EnhancedRetryRoundTripper) applyJitter(delay time.Duration) time.Duration {
	if rt.config.JitterFactor <= 0 {
		return delay
	}
	
	jitterRange := float64(delay) * rt.config.JitterFactor
	jitter := (rand.Float64() * 2 - 1) * jitterRange // Random between -jitterRange and +jitterRange
	
	newDelay := time.Duration(float64(delay) + jitter)
	if newDelay < 0 {
		return 0
	}
	
	return newDelay
}

// readRetryAfterHeader reads the Retry-After header from the response
func (rt *EnhancedRetryRoundTripper) readRetryAfterHeader(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}
	
	// Try to parse as seconds (integer)
	waitSeconds, err := time.ParseDuration(retryAfter + "s")
	if err == nil {
		return waitSeconds
	}
	
	// Try to parse as HTTP date
	retryTime, err := http.ParseTime(retryAfter)
	if err == nil {
		return time.Until(retryTime)
	}
	
	return 0
}

// cloneRequest creates a shallow copy of the request
func (rt *EnhancedRetryRoundTripper) cloneRequest(req *http.Request) *http.Request {
	reqCopy := new(http.Request)
	*reqCopy = *req
	
	// Deep copy the header
	if req.Header != nil {
		reqCopy.Header = make(http.Header, len(req.Header))
		for k, v := range req.Header {
			reqCopy.Header[k] = append([]string(nil), v...)
		}
	}
	
	return reqCopy
}

// RetryableClient wraps a client with retry logic
type RetryableClient struct {
	*Client
	retryConfig *RetryConfig
}

// NewRetryableClient creates a new client with enhanced retry capabilities
func NewRetryableClient(client *Client, config *RetryConfig) *RetryableClient {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	// Wrap the existing HTTP client with retry logic
	if client.HTTPClient != nil {
		client.HTTPClient.Transport = NewEnhancedRetryRoundTripper(
			client.HTTPClient.Transport,
			config,
		)
	}
	
	return &RetryableClient{
		Client:      client,
		retryConfig: config,
	}
}

// RetryWithContext executes a function with retry logic
func RetryWithContext(ctx context.Context, config *RetryConfig, fn func() error) error {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	var lastErr error
	
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}
		
		if attempt < config.MaxRetries {
			// Calculate delay
			delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt)))
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
			
			// Wait before retrying
			select {
			case <-time.After(delay):
				// Continue to next retry
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	
	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	// Check for specific Power BI error types that are retryable
	if httpErr, ok := err.(HTTPUnsuccessfulError); ok {
		// Rate limiting errors
		if httpErr.Response != nil && httpErr.Response.StatusCode == 429 {
			return true
		}
		// Server errors
		if httpErr.Response != nil && httpErr.Response.StatusCode >= 500 && httpErr.Response.StatusCode < 600 {
			return true
		}
	}
	
	return false
}