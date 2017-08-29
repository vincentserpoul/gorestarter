package mid

import (
	"context"
	"net/http"

	"github.com/segmentio/ksuid"
)

const (
	contextKeyRequestID = ContextKey("requestID")
)

// RequestID adds a requestID to the request context
// ksuid is a unique global id that is orderable by time (a step up normal uuid)
func RequestID() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := ksuid.New()
			w.Header().Add("requestID", requestID.String())
			ctx := context.WithValue(r.Context(), contextKeyRequestID, requestID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID will retrieve the request id from the context if there is one
func GetRequestID(ctx context.Context) ksuid.KSUID {
	if reqID, ok := ctx.Value(contextKeyRequestID).(ksuid.KSUID); ok {
		return reqID
	}

	return ksuid.Nil
}
