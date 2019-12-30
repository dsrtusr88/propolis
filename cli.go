package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/fs"
)

const (
	usage = `
	
Description:
    Make sure files are in good shape before uploading.
	Detect trumpable releases.
	
Usage:
    propolis [-c|--check-for-dupes] <PATH>

Options:
    -c, --check-for-dupes  Check the tracker for duplicates.
    -h, --help             Show this screen.
    --version              Show version.
`
	fullName    = "propolis"
	fullVersion = "%s -- v%s"
	version     = "0.1.0"
)

func userAgent() string {
	return fullName + "/" + version
}

type propolisArgs struct {
	builtin       bool
	checkForDupes bool
	path          string
}

func (m *propolisArgs) parseCLI(osArgs []string) error {
	// parse arguments and options
	args, err := docopt.Parse(fmt.Sprintf(usage, version), osArgs, true, fmt.Sprintf(fullVersion, fullName, version), false, false)
	if err != nil {
		return errors.Wrap(err, "incorrect arguments")
	}
	if len(args) == 0 {
		// builtin command, nothing to do.
		m.builtin = true
		return nil
	}
	m.checkForDupes = args["--check-for-dupes"].(bool)
	m.path = args["<PATH>"].(string)
	if !fs.DirExists(m.path) {
		return errors.New("target path " + m.path + " not found")
	}
	return nil
}
