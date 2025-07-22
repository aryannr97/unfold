package base64

import (
	"os"
	"strings"
	"testing"
)

func Test_commandDecodeConfig_Execute(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Valid base64 string",
			args: "VGhpcyBpcyBhIHRlc3Qgc3RyaW5nIGZvciB1bmZvbGQ=",
			want: "This is a test string for unfold",
		},
		{
			name: "Invalid base64 string",
			args: "invalid",
			want: "failed to decode given string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandDecodeConfig
			if c.GetFlagSet() == nil {
				t.Errorf("commandDecodeConfig.GetFlagSet() = nil, want non-nil")
			}
			os.Args = []string{"unfold", "base64", "decode", tt.args}
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandDecodeConfig.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
