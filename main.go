package main

import (
	"fmt"
	"os"

	"gitlab.com/catastrophic/assistance/logthis"
)

func main() {
	// parsing CLI
	cli := &propolisArgs{}
	if err := cli.parseCLI(os.Args[1:]); err != nil {
		logthis.Error(err, logthis.NORMAL)
		return
	}
	if cli.builtin {
		return
	}

	fmt.Println(cli.path)
	fmt.Println(cli.checkForDupes)

}
