package main

import (
	"os"
	"os/exec"

	"github.com/dotqin/blf/log"
)

var cmd *exec.Cmd

func RunCmd(path string) {

	var root = os.Getenv("GOROOT")
	cmd = exec.Command(root+"/bin/go", "run", path+"/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go cmd.Run()
}

func Kill() {
	defer func() {
		if e := recover(); e != nil {
			log.Debug("Kill.recover -> ", e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Debug("Kill -> ", err)
		}
	}
}
