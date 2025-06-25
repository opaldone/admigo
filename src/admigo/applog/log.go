package applog

import (
	"fmt"
	"os"
)

func Danger(step string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] ", step)
	fmt.Fprintln(os.Stderr, args...)
}
