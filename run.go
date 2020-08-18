package propolis

import (
	"errors"
	"path/filepath"
	"strconv"

	"gitlab.com/catastrophic/assistance/flac"
	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/music"
	"gitlab.com/catastrophic/assistance/ui"
)

var (
	log = &Log{}
)

func Run(path string, disableSpecs, problemsOnly, snatched bool, version string) (*Propolis, string, error) {
	// setting output config
	log.problemsOnly = problemsOnly

	// by default, metadata (spectrograms, etc), will be put in a side folder.
	metadataDir := path + " (Metadata)"
	if snatched {
		metadataDir = filepath.Join(path, "Metadata")
	}
	release := music.NewWithExternalMetadata(path, metadataDir)
	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)

	// creating overall check struct and adding the first checks
	res := NewPropolis(path, StdOutput)

	// general checks
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking Path is a music release"), logthis.NORMAL)
	err := release.ParseFiles()

	res.ErrorCheck(LevelCritical, "2.3.1", "Release contains FLAC files", "Error parsing files", err, DoNotAppendError)
	if err != nil {
		res.ConditionCheck(LevelCritical, "2.2.10.8", arrowHeader+"At least one FLAC has illegal ID3v2 tags.", arrowHeader+err.Error(), errors.Is(err, flac.ErrNoFlacHeader))
	}
	res.ConditionCheck(LevelInfo, internalRule, "Total size of release folder: "+strconv.FormatFloat(totalSize, 'f', 2, 32)+"Mb.", "", true)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleMusic), logthis.NORMAL)
	res = CheckMusicFiles(release, res)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleOrganization), logthis.NORMAL)
	res = CheckOrganization(release, snatched, res)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleTags), logthis.NORMAL)
	res = CheckTags(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFilenames), logthis.NORMAL)
	res = CheckFilenames(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleExtraFiles), logthis.NORMAL)
	res = CheckExtraFiles(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFoldername), logthis.NORMAL)
	res = CheckFolderName(release, res)

	var overviewFile string
	if !disableSpecs {
		logthis.Info(titleHeader+ui.BlueBoldUnderlined("Generating spectrograms"), logthis.NORMAL)
		overviewFile, err = GenerateSpectrograms(release)
		if err != nil {
			logthis.Error(err, logthis.NORMAL)
		} else {
			logthis.Info(ui.BlueBold("Spectrograms generated in "+metadataDir+". Check for transcodes (see wiki#408)."), logthis.NORMAL)
		}
	}
	logthis.Info("\n"+titleHeader+ui.BlueBoldUnderlined("Results\n")+ui.Blue(res.Summary()), logthis.NORMAL)
	// saving log to file
	if err != res.SaveOuput(metadataDir, version) {
		return res, overviewFile, err
	}
	return res, overviewFile, nil
}
