//go:build ignore

package main

import (
	"os"
	"os/exec"
	"slices"
)

func main() {
	var cmd string

	switch {
	case slices.Contains(os.Args, "--build"):
		cmd = "go build -o Wike.exe ./cmd/wike/main.go"
	case slices.Contains(os.Args, "--run"):
		cmd = "go run ./cmd/wike/main.go"
	case slices.Contains(os.Args, "--release"):
		cmd = "bunx relion"
	}

	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	execCmd.Env = append(os.Environ(), "GOOS=windows")
	execCmd.Run()
}
