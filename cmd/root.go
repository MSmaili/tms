package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tmuxctl <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		Start(os.Args[2:])
	// case "switch":
	// 	SwitchCmd(os.Args[2:])
	// case "list":
	// 	ListCmd(os.Args[2:])
	// case "load":
	// 	LoadCmd(os.Args[2:])
	case "create":
		Test(os.Args[2:])
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func Test(flags []string) {
	fmt.Println("test")
	fmt.Println(flags)
}
