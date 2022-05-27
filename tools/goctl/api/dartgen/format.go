package dartgen

import (
	"fmt"
	"os"
	"os/exec"
)

const dartExec = "dart"

func formatDir(dir string) error {
	ok, err := dirctoryExists(dir)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("format failed, directory %q does not exist", dir)
	}

	_, err = exec.LookPath(dartExec)
	if err != nil {
		return err
	}
	cmd := exec.Command(dartExec, "format", dir)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func dirctoryExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
