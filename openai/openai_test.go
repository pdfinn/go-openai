package openai

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func mockHTTPClient(mockResponse string, mockStatusCode int) *http.Client {
	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(mockStatusCode)
		w.Write([]byte(mockResponse))
	}))

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(mockServer.URL)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Add this line
		},
	}
	return httpClient
}

func TestCallOpenAI(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		mockResponse string
		mockStatus   int
		wantErr      error
		wantResponse string
	}{
		{
			name: "valid response",
			config: Config{
				APIKey:      "valid-api-key",
				Instruction: "Translate the following English text to French: '{}'.",
				Input:       "Hello, World!",
				Temperature: 0.7,
				Model:       "gpt-3.5-turbo",
			},
			mockResponse: `{"choices": [{"text": "Bonjour le monde!"}]}`,
			mockStatus:   http.StatusOK,
			wantErr:      nil,
			wantResponse: "Bonjour le monde!",
		},
		{
			name: "invalid model",
			config: Config{
				APIKey:      "valid-api-key",
				Instruction: "Translate the following English text to French: '{}'.",
				Input:       "Hello, World!",
				Temperature: 0.7,
				Model:       "invalid-model",
			},
			mockResponse: "",
			mockStatus:   http.StatusOK,
			wantErr:      ErrInvalidModel,
			wantResponse: "",
		},
		{
			name: "API error",
			config: Config{
				APIKey:      "valid-api-key",
				Instruction: "Translate the following English text to French: '{}'.",
				Input:       "Hello, World!",
				Temperature: 0.7,
				Model:       "gpt-3.5-turbo",
			},
			mockResponse: "",
			mockStatus:   http.StatusBadRequest,
			wantErr:      ErrNoChoices,
			wantResponse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := mockHTTPClient(tt.mockResponse, tt.mockStatus)
			got, err := CallOpenAI(tt.config, httpClient)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}

			if got != tt.wantResponse {
				t.Errorf("got response %s, want %s", got, tt.wantResponse)
			}
		})
	}
}

func TestValidateModel(t *testing.T) {
	tests := []struct {
		model string
		want  bool
	}{
		{"gpt-4", true},
		{"gpt-3.5-turbo", true},
		{"invalid-model", false},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			if got := ValidateModel(tt.model); got != tt.want {
				t.Errorf("ValidateModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
