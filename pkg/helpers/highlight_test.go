package helpers

import "testing"

func TestValueFormatting(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		fn   func(string) string
		want string
	}{
		{
			name: "GreenValue",
			args: args{
				val: "test",
			},
			fn:   GreenValue,
			want: "\033[1m\033[32mtest\033[0m",
		},
		{
			name: "RedValue",
			args: args{
				val: "test",
			},
			fn:   RedValue,
			want: "\033[1m\033[31mtest\033[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.args.val); got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
