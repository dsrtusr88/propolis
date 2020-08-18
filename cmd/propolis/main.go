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

	res, _, err := propolis.Run(cli.path, cli.disableSpecs, cli.problemsOnly, cli.snatched, Version)
	if err != nil {
		logthis.Error(err, logthis.NORMAL)
	}

	// returning nonzero exit status if something serious was found
	if res.CriticalErrors != 0 {
		syscall.Exit(1)
	}
}
