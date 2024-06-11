package system

import (
	"bytes"
	"fmt"
	"lance-light/internal/log"
	"os/exec"
	"strings"
)

func ExecCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stderr, err := cmd.CombinedOutput()

	if err == nil {
		log.MsgFatalAndExit(err, fmt.Sprintf("Failed to exec.　| \"%s\" >>> %s",
			strings.Join(cmd.Args, " "),
			string(stderr)))
	}
}

func ExecCommandGetResult(command string, args []string) []string {
	cmd := exec.Command(command, args...)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	// 実行
	err := cmd.Run()
	if err != nil {
		log.MsgFatalAndExit(err, fmt.Sprintf("Failed to exec.　| \"%s\" >>> %s",
			strings.Join(cmd.Args, " "),
			stderr.String()))
	}

	output := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	return output
}
