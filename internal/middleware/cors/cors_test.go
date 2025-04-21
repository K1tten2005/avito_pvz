package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		method            string
		expectedHeaders   map[string]string
		expectedStatusCode int
	}{
		{
			name:   "GET request",
			method: http.MethodGet,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods":    "POST,GET",
				"Access-Control-Allow-Headers":    "Authorization,Content-Type,X-Csrf-Token",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":   "Authorization,X-Csrf-Token",
				"Access-Control-Allow-Origin":     "http://localhost:3000",
				"Access-Control-Max-Age":          "86400",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "OPTIONS request",
			method: http.MethodOptions,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods":    "POST,GET",
				"Access-Control-Allow-Headers":    "Authorization,Content-Type,X-Csrf-Token",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":   "Authorization,X-Csrf-Token",
				"Access-Control-Allow-Origin":     "http://localhost:3000",
				"Access-Control-Max-Age":          "86400",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := CorsMiddleware(nextHandler)

			req := httptest.NewRequest(test.method, "http://localhost", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			for key, expectedValue := range test.expectedHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(key), "Заголовок %s должен быть %s", key, expectedValue)
			}
		})
	}
}
