package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDisableCache(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	disabledCacheHandler := DisableCache(handler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	disabledCacheHandler.ServeHTTP(rr, req)

	if cacheControl := rr.Header().Get("Cache-Control"); cacheControl != "no-cache, no-store, must-revalidate" {
		t.Errorf("handler returned wrong Cache-Control header: got %v want %v", cacheControl, "no-cache, no-store, must-revalidate")
	}
}
