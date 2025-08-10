package helpers

import (
	"testing"
)

func TestGracefullyExit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test gracefully exit with panic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer GracefullyExit()
				panic("test panic")
			}()
		})
	}
}
