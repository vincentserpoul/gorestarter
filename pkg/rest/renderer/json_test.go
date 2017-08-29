package renderer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseJSONRender(t *testing.T) {
	r, _ := http.NewRequest("GET", ``, nil)
	type args struct {
		w http.ResponseWriter
		r *http.Request
		e interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "working json marshalling",
			args: args{
				w: httptest.NewRecorder(),
				r: r,
				e: ErrResponse{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResponseJSONRender(tt.args.w, tt.args.r, tt.args.e)
		})
	}
}
