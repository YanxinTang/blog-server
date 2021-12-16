package cmd

import (
	"fmt"
	"os"
)

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}

func exitWithError(err error) {
	exitf("Exit: %s", err)
	os.Exit(1)
}
