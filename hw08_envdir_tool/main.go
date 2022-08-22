package main

import (
	"os"
)

func main() {
	args := os.Args
	env, err := ReadDir(args[1])
	if err != nil {
		return
	}
	RunCmd(args[2:], env)
}
