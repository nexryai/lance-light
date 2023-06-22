package core

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stderr, err := cmd.CombinedOutput()

	ExitOnError(err, fmt.Sprintf("Failed to exec.　| \"%s\" >>> %s",
		strings.Join(cmd.Args, " "),
		string(stderr)))

}
