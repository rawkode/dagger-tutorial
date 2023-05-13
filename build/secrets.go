package devtool

import (
	"os/exec"
	"strings"
)

func GetSecret(lookupRef string) (string, error) {
	opGetSecret := exec.Command("op", "read", lookupRef)
	stdout, err := opGetSecret.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(stdout), "\n"), nil
}
