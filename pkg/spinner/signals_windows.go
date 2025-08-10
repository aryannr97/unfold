//go:build windows

package spinner

import "os"

// Windows lacks many UNIX signals; only handle Ctrl+C equivalent.
var handledSignals = []os.Signal{os.Interrupt}
