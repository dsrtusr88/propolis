package main

import (
	"os"

	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/music"
)

const (
	internalRule = "internal"
)

var (
	log = &Log{}
)

func main() {
	// checking external tools
	if err := checkExternalBinaries("sox", "flac", "metaflac"); err != nil {
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

	release := music.New(cli.path)

	logthis.Info("Checking Path is a music release", logthis.NORMAL)
	err := release.ParseFiles()
	log.CriticalResult(err == nil, internalRule, "Release contains FLAC files", "Error parsing files")
	if err != nil {
		log.BadResult(err == nil, internalRule, "", "Critical error: "+err.Error())
		return
	}

	logthis.Info("Checking music files", logthis.NORMAL)
	if err := CheckMusicFiles(release); err != nil {
		log.BadResult(err == nil, internalRule, "", "Critical error: "+err.Error())
		return
	}

	logthis.Info("Checking organization", logthis.NORMAL)
	if err := CheckOrganization(release); err != nil {
		log.BadResult(err == nil, internalRule, "", "Critical error: "+err.Error())
		return
	}
	// in folder or subfolders for CDs

	logthis.Info("Checking tags", logthis.NORMAL)
	if err := CheckTags(release); err != nil {
		log.BadResult(err == nil, internalRule, "", "Critical error: "+err.Error())
		return
	}

	logthis.Info("Checking filenames", logthis.NORMAL)
	// 2.3.13.filenames contain track
	// filenames contain at least beginning of title
	// 2.3.14. check numbers at beginning of filenames
	// 2.3.15. no 2 same numbers in folder
	// flac is small caps
	// 2.3.20. Leading spaces are not allowed in any file or folder names + leading dots
	logthis.Info("Checking extra files", logthis.NORMAL)
	// check size of side art + % of total size
	// check forbidden extensions
	logthis.Info("Checking folder name", logthis.NORMAL)

	logthis.Info("Generating spectrograms", logthis.NORMAL)

}
