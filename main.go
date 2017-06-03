package main

import (
	"fmt"
	"os"

	"github.com/damselem/autohosts/cmd"
)

func main() {
	if err := cmd.GetRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
