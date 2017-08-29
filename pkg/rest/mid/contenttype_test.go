package mid

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeader(t *testing.T) {
	key := "test"
	value := "test value"
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	midWared := Header(key, value)(fakeHandler)
	rr := httptest.NewRecorder()
	request, errR := http.NewRequest("GET", ``, nil)
	if errR != nil {
		t.Fatalf("request creation failed %v", errR)
	}

	midWared.ServeHTTP(rr, request)
	if rr.Header().Get(key) != value {
		t.Errorf("expected %s, got %s", value, rr.Header().Get(key))
	}
}
