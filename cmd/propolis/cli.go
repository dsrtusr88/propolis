package main

import (
	"fmt"
	"os"
	"path/filepath"

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
    propolis [--metadata-root=<METADATA_PATH>] [--no-specs] [--no-overview] [--only-problems] [--snatched] [--json] <PATH>

Options:
    --snatched                       Snatched mode: allow varroa metadata files, spec generated in <PATH>
    --no-specs                       Disable spectrograms generation.
    --no-overview                    Disable spectrograms overview.
    --only-problems                  Only show problems (warnings & errors).
    --json                           Toggles JSON output. Sets --only-problems to false.
    --metadata-root=<METADATA_PATH>  Save propolis metadata inside this folder.
    -h, --help                       Show this screen.
    --version                        Show version.
`
	fullName    = "propolis"
	fullVersion = "%s -- %s"
)

var Version = "dev"

type propolisArgs struct {
	builtin              bool
	disableSpecs         bool
	disableCombinedSpecs bool
	problemsOnly         bool
	snatched             bool
	jsonOutput           bool
	path                 string
	metadataRoot         string
}

func (m *propolisArgs) parseCLI(osArgs []string) error {
	// parse arguments and options
	args, err := docopt.ParseArgs(fmt.Sprintf(usage, Version), osArgs, fmt.Sprintf(fullVersion, fullName, Version))
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
	m.disableCombinedSpecs = args["--no-overview"].(bool)
	m.problemsOnly = args["--only-problems"].(bool)
	m.jsonOutput = args["--json"].(bool)
	if m.jsonOutput {
		m.problemsOnly = false
	}
	m.path = filepath.Clean(args["<PATH>"].(string))
	if !fs.DirExists(m.path) {
		return errors.New("target path " + m.path + " not found")
	}
	// if given current directory, going back up to find the current directory name
	if m.path == "." {
		cwd, _ := os.Getwd()
		m.path = filepath.Join("..", filepath.Base(cwd))
	}

	metadataRoot, err := args.String("--metadata-root")
	if err == nil {
		m.metadataRoot = filepath.Clean(metadataRoot)
		if m.metadataRoot != "" && !fs.DirExists(m.metadataRoot) {
			return errors.New("target metadata root path " + m.metadataRoot + " not found")
		}
	}

	if m.snatched && m.metadataRoot != "" {
		return errors.New("--snatched implies metadata will be saved inside the release folder, not compatible with --metadata-root")
	}
	return nil
}
