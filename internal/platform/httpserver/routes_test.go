package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type readinessFunc func(context.Context) error

func (fn readinessFunc) Ping(ctx context.Context) error {
	return fn(ctx)
}

func TestReadyHandlerWithoutDatabase(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	readyHandler(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), `"database":"not_configured"`) {
		t.Fatalf("body = %s, want database not_configured", rec.Body.String())
	}
}

func TestReadyHandlerWithUnavailableDatabase(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	checker := readinessFunc(func(context.Context) error {
		return errors.New("database unavailable")
	})

	readyHandler(checker).ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
