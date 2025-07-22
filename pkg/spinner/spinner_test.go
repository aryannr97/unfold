package spinner

import (
	"reflect"
	"syscall"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	type args struct {
		spinnerType Type
	}
	tests := []struct {
		name string
		args args
		want *Spinner
	}{
		{
			name: "BrailDot spinner",
			args: args{
				spinnerType: BrailDot,
			},
			want: &Spinner{
				frames: brailDotFrames,
			},
		},
		{
			name: "ClassicDot spinner",
			args: args{
				spinnerType: ClassicDot,
			},
			want: &Spinner{
				frames: classicDotFrames,
			},
		},
		{
			name: "DualBall spinner",
			args: args{
				spinnerType: DualBall,
			},
			want: &Spinner{
				frames: dualBallFrames,
			},
		},
		{
			name: "Circle spinner",
			args: args{
				spinnerType: Circle,
			},
			want: &Spinner{
				frames: circleFrames,
			},
		},
		{
			name: "Classic spinner",
			args: args{
				spinnerType: "Default",
			},
			want: &Spinner{
				frames: classicFrames,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.spinnerType); !reflect.DeepEqual(got.frames, tt.want.frames) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Spinner(t *testing.T) {
	tests := []struct {
		name    string
		runTime func(s *Spinner)
	}{
		{
			name: "Run the spinner",
			runTime: func(s *Spinner) {
				time.Sleep(1 * time.Second)
				s.Clear()
			},
		},
		{
			name: "Run the spinner and stop on os signal",
			runTime: func(s *Spinner) {
				time.Sleep(1 * time.Second)
				osCue <- syscall.SIGINT
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Get(BrailDot)
			s.test = true
			go s.Start()
			tt.runTime(s)
			time.Sleep(1 * time.Second)
		})
	}
}
