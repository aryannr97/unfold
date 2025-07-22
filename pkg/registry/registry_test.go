package registry

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Registry
	}{
		{
			name: "Fetch the registry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); got == nil {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
