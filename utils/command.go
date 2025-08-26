package utils

import (
	"os/exec"
	"runtime"
)

func ExecuteCommand(c string) ([]byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", c)
	} else {
		cmd = exec.Command("/bin/bash", "-c", c)
	}
	return cmd.CombinedOutput()
}
