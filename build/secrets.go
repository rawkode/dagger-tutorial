package devtool

import (
	"os/exec"
)

func GetSecret(lookupRef string) (string, error) {
	opGetSecret := exec.Command("op", "read", lookupRef)
	stdout, err := opGetSecret.Output()

	if err != nil {
		return "", err
	}

	return string(stdout), nil
}
