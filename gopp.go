package main

import (
	"fmt"
	"github.com/gizak/gopp/cmd"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
