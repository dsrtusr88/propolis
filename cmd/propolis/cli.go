package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/fs"
)

const (
	usage = `
    ___  ____ ____ ___  ____ _    _ ____ 
    |__] |__/ |  | |__] |  | |    | [__  
    |    |  \ |__| |    |__| |___ | ___]  (%s)
	
Description:
    Make sure files are in good shape before uploading.
    Detect trumpable releases.
	
Usage:
    propolis [--no-specs] [--only-problems] [--snatched] [--json] <PATH>

Options:
    --snatched       Snatched mode: allow varroa metadata files, spec generated in <PATH>
    --no-specs       Disable spectrograms generation.
    --only-problems  Only show problems (warnings & errors).
    -h, --help       Show this screen.
    --version        Show version.
`
	fullName    = "propolis"
	fullVersion = "%s -- %s"
)

var Version = "dev"

type propolisArgs struct {
	builtin      bool
	disableSpecs bool
	problemsOnly bool
	snatched     bool
	jsonOutput   bool
	path         string
}

func (m *propolisArgs) parseCLI(osArgs []string) error {
	// parse arguments and options
	args, err := docopt.Parse(fmt.Sprintf(usage, Version), osArgs, true, fmt.Sprintf(fullVersion, fullName, Version), false, false)
	if err != nil {
		return errors.Wrap(err, "incorrect arguments")
	}
	if len(args) == 0 {
		// builtin command, nothing to do.
		m.builtin = true
		return nil
	}
	m.snatched = args["--snatched"].(bool)
	m.disableSpecs = args["--no-specs"].(bool)
	m.problemsOnly = args["--only-problems"].(bool)
	m.jsonOutput = args["--json"].(bool)
	m.path = args["<PATH>"].(string)
	if !fs.DirExists(m.path) {
		return errors.New("target path " + m.path + " not found")
	}
	return nil
}
