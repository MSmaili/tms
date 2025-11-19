package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunInteractive(args ...string) error {
	cmd := exec.Command("tmux", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Run(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("tmux error:", string(out))
	}

	return strings.TrimSpace(string(out)), err
}

func Attach(session string) error {
	if session == "" {
		return RunInteractive("attach-session")
	}

	return RunInteractive("attach-session", "-t", session)
}
