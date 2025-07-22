package helpers

import "fmt"

// reset is the default color
const reset = "\033[0m"

// GreenValue returns the green value
func GreenValue(val string) string {
	boldGreen := "\033[1m\033[32m"
	return fmt.Sprintf("%s%s%s", boldGreen, val, reset)
}

// RedValue returns the red value
func RedValue(val string) string {
	boldRed := "\033[1m\033[31m"
	return fmt.Sprintf("%s%s%s", boldRed, val, reset)
}
