package main

import (
	"os"
	"syscall"

	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/passelecasque/propolis"
)

func main() {
	// checking external tools
	if err := propolis.CheckExternalBinaries("sox", "flac"); err != nil {
		logthis.Error(err, logthis.NORMAL)
		return
	}

	// parsing CLI
	cli := &propolisArgs{}
	if err := cli.parseCLI(os.Args[1:]); err != nil {
		logthis.Error(err, logthis.NORMAL)
		return
	}
	if cli.builtin {
		return
	}

	results, _, err := propolis.Run(cli.path, cli.disableSpecs, cli.problemsOnly, cli.snatched, cli.jsonOutput, Version)
	if err != nil {
		logthis.Error(err, logthis.NORMAL)
	}

	// returning nonzero exit status if something serious was found
	if results.Errors != 0 {
		syscall.Exit(1)
	}
}
