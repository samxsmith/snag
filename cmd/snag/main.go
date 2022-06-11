package main

import (
	"fmt"
	"os"

	"github.com/samxsmith/snag"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Can't read working directory: %s \n", err)
	}
	snag.Run(cwd)
}
