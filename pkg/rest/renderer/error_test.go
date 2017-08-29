package renderer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/render"
)

func TestErrResponse_Render(t *testing.T) {

	tests := []struct {
		name        string
		errResponse *ErrResponse
		wantErr     bool
	}{
		{
			name:        "error response",
			errResponse: &ErrResponse{},
			wantErr:     false,
		},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", ``, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.errResponse.Render(w, r); (err != nil) != tt.wantErr {
				t.Errorf("ErrResponse.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrFuncs(t *testing.T) {

	tests := []struct {
		name      string
		funcToUse func(error) render.Renderer
		err       error
		want      render.Renderer
	}{
		{
			name:      "working error renderer",
			funcToUse: ErrRender,
			err:       errors.New("test"),
			want: &ErrResponse{
				Err:            errors.New("test"),
				HTTPStatusCode: http.StatusInternalServerError,
				StatusText:     "Error rendering response.",
			},
		},
		{
			name:      "working error renderer invalid request",
			funcToUse: ErrInvalidRequest,
			err:       errors.New("test err"),
			want: &ErrResponse{
				Err:            errors.New("test err"),
				HTTPStatusCode: http.StatusBadRequest,
				StatusText:     "Invalid request.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.funcToUse(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrRender() = %v, want %v", got, tt.want)
			}
		})
	}
}
