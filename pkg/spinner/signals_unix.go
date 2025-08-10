//go:build !windows

package spinner

import (
	"os"
	"syscall"
)

var handledSignals = []os.Signal{syscall.SIGINT, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGHUP}
