package apigen

import (
	"os/exec"
	"strings"
)

func getGitName() string {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func getGitEmail() string {
	cmd := exec.Command("git", "config", "user.email")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
