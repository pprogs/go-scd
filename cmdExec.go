package main

import (
	"bytes"
	"os/exec"
)

func execCommand(c *command) int {

	lg.Printf("Exec: %v With %v\n", c.CmdName, c.CmdArgs)

	cmd := exec.Command(c.CmdName, c.CmdArgs...)
	cmd.Dir = c.WorkDir

	var stdOut, stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	lg.Printf("stdErr: %s\n", stdErr.String())
	lg.Printf("stdOut: %s\n", stdOut.String())

	if ee, ok := err.(*exec.ExitError); ok {
		lg.Printf("exec exit ERROR: %v\n", ee)
		return ee.ExitCode()
	}

	lg.Printf("exec OK !!\n")

	return 0
}
