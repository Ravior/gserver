package gproc

import "testing"

func Test_ShellExec(t *testing.T) {
	if err := ShellExec("ls -alh"); err != nil {
		t.Fail()
	}
}
