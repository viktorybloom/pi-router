package system

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func MustRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("must run as root; use sudo")
	}
	return nil
}

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func Exists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
