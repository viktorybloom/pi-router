package system

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return string(out), err
}

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func MustRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("must run as root; use sudo")
	}

	return nil
}
