package csp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCspMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		expectedStatus   int
		expectedCspValue string
	}{
		{
			name:             "GET request with CSP header",
			method:           http.MethodGet,
			expectedStatus:   http.StatusOK,
			expectedCspValue: "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			rr := httptest.NewRecorder()

			handler := CspMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			cspHeader := rr.Header().Get("Content-Security-Policy")
			if tt.expectedCspValue != "" {
				assert.Equal(t, tt.expectedCspValue, cspHeader)
			} else {
				assert.Empty(t, cspHeader)
			}
		})
	}
}
