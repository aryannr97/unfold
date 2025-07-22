package spinner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Type represents the type of spinner
type Type string

const (
	BrailDot   Type = "brailDot"
	ClassicDot Type = "classicDot"
	DualBall   Type = "dualBall"
	Circle     Type = "circle"
)

var (
	// osCue is a channel to listen for os signals
	osCue            = make(chan os.Signal, 1)
	brailDotFrames   = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	classicDotFrames = []string{".  ", ".. ", "..."}
	classicFrames    = []string{"|", "/", "-", "\\"}
	dualBallFrames   = []string{" ·", "● ", "●●", " ●", "· "}
	circleFrames     = []string{"◐", "◓", "◑", "◒"}
)

// Spinner represents the cli spinner
type Spinner struct {
	frames []string
	stop   chan bool
	done   chan bool
	test   bool
}

// Start starts the spinner
func (s *Spinner) Start() {
	fmt.Print("\033[?25l") // Hide cursor
	s.run()
}

// Clear stops the spinner
func (s *Spinner) Clear() {
	s.stop <- true
	<-s.done
}

// Get returns the spinner
func Get(spinnerType Type) *Spinner {
	var frames []string
	switch spinnerType {
	case BrailDot:
		frames = brailDotFrames
	case ClassicDot:
		frames = classicDotFrames
	case DualBall:
		frames = dualBallFrames
	case Circle:
		frames = circleFrames
	default:
		frames = classicFrames
	}
	return &Spinner{
		frames: frames,
		stop:   make(chan bool, 1),
		done:   make(chan bool),
	}
}

// run renders the spinner
func (s *Spinner) run() {
	signal.Notify(osCue, syscall.SIGINT, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGHUP)
	i := 0
	for {
		select {
		case <-osCue:
			fmt.Println("Clearing spinner")
			ShowCursor()
			// If test is true, don't exit the program as error
			if !s.test {
				os.Exit(1)
			}
		case <-s.stop:
			fmt.Printf("\r \r")
			ShowCursor()
			s.done <- true
			return
		default:
			fmt.Printf("\r%s", s.frames[i%len(s.frames)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ShowCursor shows the cursor
func ShowCursor() {
	fmt.Print("\033[?25h") // Show cursor
}
