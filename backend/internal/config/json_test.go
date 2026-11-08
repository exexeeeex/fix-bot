package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	recorder := httptest.NewRecorder()

	WriteError(recorder, http.StatusBadRequest, "invalid request")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want %q", contentType, "application/json; charset=utf-8")
	}
	if body := recorder.Body.String(); body != "{\"error\":\"invalid request\"}\n" {
		t.Fatalf("body = %q, want %q", body, "{\"error\":\"invalid request\"}\n")
	}
}
