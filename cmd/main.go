package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/m4tthewde/huego/pkg/frontend"
)

func main() {
	log.SetOutput(io.Discard)

	p := frontend.NewProgram()
	m, err := p.Run()
	if err != nil {
		fmt.Printf("frontend err: %v\n", err)
		os.Exit(1)
	}

	model := m.(frontend.Model)
	if model.Err != nil {
		fmt.Printf("backend err: %v\n", model.Err)
		os.Exit(1)
	}

}
