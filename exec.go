package main

import (
	"fmt"
	"os/exec"
)

func execCmd(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		fmt.Println("execCmd", name, arg, "returned:", err, "; output: ", out)
		return ""
	}
	return string(out)
}
