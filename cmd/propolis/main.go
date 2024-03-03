package main

import (
    "os"
    "syscall"

    "gitlab.com/catastrophic/assistance/logthis"
    "gitlab.com/passelecasque/propolis"
)

type propolisArgs struct {
    // Define the fields of propolisArgs struct here
}

// Define Version constant
const Version = "1.0.0"

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

    results, _, err := propolis.Run(cli.path, cli.metadataRoot, cli.disableSpecs, cli.disableCombinedSpecs, cli.problemsOnly, cli.snatched, cli.jsonOutput, true, Version)
    if err != nil {
        logthis.Error(err, logthis.NORMAL)
    }

    // returning nonzero exit status if something serious was found
    if results.Errors != 0 {
        syscall.Exit(1)
    }
}
