package main

import (
	"os/exec"

	"github.com/pkg/errors"
)

func checkExternalBinaries(externalBinaries ...string) error {
	// check the required binaries are installed
	for _, r := range externalBinaries {
		_, err := exec.LookPath(r)
		if err != nil {
			return errors.New(r + " is not available on this system")
		}
	}
	return nil
}
