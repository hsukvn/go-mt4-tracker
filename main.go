package main

import (
	"os"

	"github.com/hsukvn/go-mt4-tracker/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
