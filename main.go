package main

import (
	"os"
	"strconv"

	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/music"
	"gitlab.com/catastrophic/assistance/ui"
)

const (
	internalRule = " -- "
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

	// by default, metadata (spectrograms, etc), will be put in a side folder.
	metadataDir := cli.path + " (Metadata)"
	release := music.NewWithExternalMetadata(cli.path, metadataDir)
	res := &Results{}
	var err error

	logthis.Info("Checking Path is a music release", logthis.NORMAL)
	err = release.ParseFiles()
	res.Add(log.CriticalResult(err == nil, "2.3.1", "Release contains FLAC files", "Error parsing files"))
	if err != nil {
		res.Add(log.BadResultInfo(err == nil, "2.3.1", "", "⮕ Critical error: "+err.Error()))
		return
	}
	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)
	res.Add(log.NeutralResult(true, internalRule, "Total size of release folder: "+strconv.FormatFloat(totalSize, 'f', 2, 32)+"Mb.", ""))

	logthis.Info("Checking music files", logthis.NORMAL)
	res = CheckMusicFiles(release, res)
	logthis.Info("Checking organization", logthis.NORMAL)
	res = CheckOrganization(release, res)
	logthis.Info("Checking tags", logthis.NORMAL)
	res = CheckTags(release, res)
	logthis.Info("Checking filenames", logthis.NORMAL)
	res = CheckFilenames(release, res)
	logthis.Info("Checking extra files", logthis.NORMAL)
	res = CheckExtraFiles(release, res)
	logthis.Info("Checking folder name", logthis.NORMAL)
	res = CheckFolderName(release, res)
	logthis.Info("Generating spectrograms", logthis.NORMAL)
	if _, err := GenerateSpectrograms(release); err != nil {
		logthis.Error(err, logthis.NORMAL)
	}
	logthis.Info(ui.BlueBold("Spectrograms generated in "+metadataDir), logthis.NORMAL)
	logthis.Info(ui.BlueBoldUnderlined("\nResults:\n")+ui.Blue("⮕ "+res.String()), logthis.NORMAL)
}
