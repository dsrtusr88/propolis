package main

import (
	"os"
	"strconv"
	"syscall"

	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/music"
	"gitlab.com/catastrophic/assistance/ui"
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

	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking Path is a music release"), logthis.NORMAL)
	err = release.ParseFiles()
	res.Add(log.CriticalResult(err == nil, "2.3.1", "Release contains FLAC files", "Error parsing files"))
	if err != nil {
		res.Add(log.BadResultInfo(err == nil, "2.3.1", "", arrowHeader+"Critical error: "+err.Error()))
		return
	}
	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)
	res.Add(log.NeutralResult(true, internalRule, "Total size of release folder: "+strconv.FormatFloat(totalSize, 'f', 2, 32)+"Mb.", ""))

	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking music files"), logthis.NORMAL)
	res = CheckMusicFiles(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking organization"), logthis.NORMAL)
	res = CheckOrganization(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking tags"), logthis.NORMAL)
	res = CheckTags(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking filenames"), logthis.NORMAL)
	res = CheckFilenames(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking extra files"), logthis.NORMAL)
	res = CheckExtraFiles(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking folder name"), logthis.NORMAL)
	res = CheckFolderName(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Generating spectrograms"), logthis.NORMAL)
	if err := GenerateSpectrograms(release); err != nil {
		logthis.Error(err, logthis.NORMAL)
	} else {
		logthis.Info(ui.BlueBold("Spectrograms generated in "+metadataDir+". Check for transcodes (see wiki#408)."), logthis.NORMAL)
	}
	logthis.Info("\n"+titleHeader+ui.BlueBoldUnderlined("Results\n")+ui.Blue(res.String()), logthis.NORMAL)

	// returning nonzero exit status if something serious was found
	if res.ko != 0 {
		syscall.Exit(1)
	}
}
