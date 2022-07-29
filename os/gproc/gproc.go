package gproc

import (
	"fmt"
	"os/exec"
)

func ShellExec(command string) error {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	return err
}
