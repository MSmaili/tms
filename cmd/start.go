package cmd

import (
	"fmt"
	"os"

	"github.com/MSmaili/tmx/internal/tmux"
)

func Start(args []string) {

	var session string
	if len(args) > 0 {
		session = args[0]
	}

	err := tmux.Attach(session)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
