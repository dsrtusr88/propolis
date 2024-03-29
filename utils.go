package propolis

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func CheckExternalBinaries(externalBinaries ...string) error {
	// check the required binaries are installed
	for _, r := range externalBinaries {
		_, err := exec.LookPath(r)
		if err != nil {
			return errors.New(r + " is not available on this system")
		}
	}
	return nil
}

func IgnoreVarroaFiles(files []string) []string {
	var clean []string
	for _, e := range files {
		if !strings.Contains(e, "TrackerMetadata") && !strings.Contains(e, "Metadata") {
			clean = append(clean, e)
		}
	}
	return clean
}

func simplifyBullet(bullet string) string {
	switch bullet {
	case OKString:
		return "OK"
	case KOString:
		return "KO"
	case WarningString:
		return "!!"
	case NeutralString:
		return "--"
	default:
		return "  "
	}
}
