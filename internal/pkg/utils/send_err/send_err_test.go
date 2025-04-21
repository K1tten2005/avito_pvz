package send_err

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		errorMessage string
		expectedBody map[string]string
	}{
		{
			name:         "Test 400 Bad Request",
			statusCode:   400,
			errorMessage: "Bad Request",
			expectedBody: map[string]string{"message": "Bad Request"},
		},
		{
			name:         "Test 404 Not Found",
			statusCode:   404,
			errorMessage: "Not Found",
			expectedBody: map[string]string{"message": "Not Found"},
		},
		{
			name:         "Test 500 Internal Server Error",
			statusCode:   500,
			errorMessage: "Internal Server Error",
			expectedBody: map[string]string{"message": "Internal Server Error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			SendError(recorder, tt.errorMessage, tt.statusCode)

			assert.Equal(t, tt.statusCode, recorder.Code)

			var responseBody map[string]string
			err := json.NewDecoder(recorder.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("Error decoding response body: %v", err)
			}
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
