package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// infoLog.Printf("run command %s with args %v", cmd[0], cmd[1:])
	comm := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	comm.Stdin, comm.Stdout, comm.Stderr = os.Stdin, os.Stdout, os.Stderr

	envSl := make([]string, 0, len(env))
	for key, val := range env {
		envSl = append(envSl, fmt.Sprintf("%v=%s", key, val.Value))
	}
	// infoLog.Printf("run command with env: %v", envSl)
	comm.Env = append(os.Environ(), envSl...)

	if err := comm.Run(); err != nil {
		log.Fatal(err)
	}

	return comm.ProcessState.ExitCode()
}
