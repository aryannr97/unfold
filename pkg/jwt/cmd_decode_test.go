package jwt

import (
	"os"
	"strings"
	"testing"
)

func Test_commandDecodeConfig_Execute(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "Valid token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			want:  `unfolded JWT Claims (JSON)`,
		},
		{
			name:  "Invalid token",
			token: "invalid",
			want:  `failed to parse token`,
		},
		{
			name:  "Token with invalid claims type",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.bnVsbA.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", // payload is "null"
			want:  `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandDecodeConfig
			if c.GetFlagSet() == nil {
				t.Errorf("commandDecodeConfig.GetFlagSet() = nil, want non-nil")
			}
			os.Args = []string{"unfold", "jwt", "decode", tt.token}
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandDecodeConfig.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
