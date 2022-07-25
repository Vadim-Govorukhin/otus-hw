package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cm := exec.Command(cmd[0], cmd[1:]...)
	cm.Stdin, cm.Stdout, cm.Stderr = os.Stdin, os.Stdout, os.Stderr
	infoLog.Printf("run command %s with args %v", cmd[0], cmd[1:])

	env_sl := make([]string, 0, len(env))
	for key, val := range env {
		if !val.NeedRemove {
			env_sl = append(env_sl, fmt.Sprintf("%v=%s", key, val.Value))
		}
	}
	infoLog.Printf("run command with env: %v", env_sl)
	cm.Env = append(os.Environ(), env_sl...)

	if err := cm.Run(); err != nil {
		errorLog.Fatal(err)
	}

	return
}
