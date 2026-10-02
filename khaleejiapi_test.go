package khaleejiapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalizedAPIErrorMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{
				"code":      "MISSING_PARAMETER",
				"message":   "email is required",
				"messageEn": "email is required",
				"messageAr": "حقل email مطلوب",
			},
		})
	}))
	defer server.Close()

	client := NewWithConfig(Config{
		APIKey:     "kapi_live_test",
		BaseURL:    server.URL,
		MaxRetries: 0,
		HTTPClient: server.Client(),
	})

	_, err := client.Validation.ValidateEmail(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.MessageEn != "email is required" {
		t.Fatalf("expected english message, got %q", apiErr.MessageEn)
	}

	if apiErr.MessageAr != "حقل email مطلوب" {
		t.Fatalf("expected arabic message, got %q", apiErr.MessageAr)
	}

	if apiErr.LocalizedMessage("ar") != "حقل email مطلوب" {
		t.Fatalf("expected arabic localized message, got %q", apiErr.LocalizedMessage("ar"))
	}
}
