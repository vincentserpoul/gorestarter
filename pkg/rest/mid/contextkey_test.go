package mid

import "testing"

func TestContextKey_String(t *testing.T) {
	tests := []struct {
		name string
		c    ContextKey
		want string
	}{
		{
			name: "contextkey string ok",
			c:    ContextKey("context key"),
			want: "rest request context key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ContextKey.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
