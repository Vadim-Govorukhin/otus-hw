package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cm := exec.Command(cmd[0], cmd[1:]...)
	cm.Stdin, cm.Stdout, cm.Stderr = os.Stdin, os.Stdout, os.Stderr
	infoLog.Printf("run command %s with args %v", cmd[0], cmd[1:])

	envSl := make([]string, 0, len(env))
	for key, val := range env {
		if !val.NeedRemove {
			envSl = append(envSl, fmt.Sprintf("%v=%s", key, val.Value))
		}
	}
	infoLog.Printf("run command with env: %v", envSl)
	cm.Env = append(os.Environ(), envSl...)

	if err := cm.Run(); err != nil {
		log.Fatal(err)
	}

	return
}
