package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()

	handler := rl.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Helper to make requests
	makeRequest := func(ip string) int {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w.Code
	}

	// Test case 1: Allow up to 10 requests (burst size) from same IP with DIFFERENT ports
	// This verifies that we are rate limiting by IP, not by IP:Port
	ipBase := "192.168.1.1"
	for i := 0; i < 10; i++ {
		// Use different port for each request
		remoteAddr := fmt.Sprintf("%s:%d", ipBase, 1234+i)
		code := makeRequest(remoteAddr)
		assert.Equal(t, http.StatusOK, code, "Request %d from %s should succeed", i+1, remoteAddr)
	}

	// Test case 2: The 11th request should fail, even with a new port
	remoteAddr := fmt.Sprintf("%s:%d", ipBase, 9999)
	code := makeRequest(remoteAddr)
	assert.Equal(t, http.StatusTooManyRequests, code, "Request 11 from %s should fail", remoteAddr)

	// Test case 3: Different IP should be allowed
	ip2 := "192.168.1.2:1234"
	code = makeRequest(ip2)
	assert.Equal(t, http.StatusOK, code, "Request from new IP should succeed")
}
