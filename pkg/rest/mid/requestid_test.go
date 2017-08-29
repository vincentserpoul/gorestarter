package mid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/ksuid"
)

func TestRequestID(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	midWared := RequestID()(fakeHandler)
	rr := httptest.NewRecorder()
	request, errR := http.NewRequest("GET", ``, nil)
	if errR != nil {
		t.Fatalf("request creation failed %v", errR)
	}

	midWared.ServeHTTP(rr, request)
	if rr.Header().Get("requestID") == "" {
		t.Errorf("expected %s, got nothing", rr.Header().Get("requestID"))
	}
}

func TestGetRequestID(t *testing.T) {
	var reqID ksuid.KSUID

	ctx := context.Background()
	reqID = GetRequestID(ctx)
	if reqID != ksuid.Nil {
		t.Errorf("expected nothing, got %s", reqID)
	}

	ctx = context.WithValue(ctx, contextKeyRequestID, "testString")
	reqID = GetRequestID(ctx)
	if reqID != ksuid.Nil {
		t.Errorf("expected nothing, got %s", reqID)
	}

	requestID := ksuid.New()
	ctx = context.WithValue(ctx, contextKeyRequestID, requestID)
	reqID = GetRequestID(ctx)
	if reqID != requestID {
		t.Errorf("expected %s, got %s", requestID.String(), reqID)
	}
}
