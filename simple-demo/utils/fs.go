package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ExecPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	_, err = filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	var p = ""
	p, err = filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	return p
}
