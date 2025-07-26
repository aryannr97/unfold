package helpers

import (
	"testing"
)

func TestGetErrorResponseBody(t *testing.T) {
	type args struct {
		res []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test error response body",
			args: args{
				res: []byte(`{"error": "test error"}`),
			},
			want: "{\n  \"error\": \"test error\"\n}",
		},
		{
			name: "test error response body with multiple fields",
			args: args{
				res: []byte(`{"error": "test error", "code": 400}`),
			},
			want: "{\n  \"error\": \"test error\",\n  \"code\": 400\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorResponseBody(tt.args.res); len(got) == len(tt.want) {
				t.Errorf("GetErrorResponseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
