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

func Run(path string, disableSpecs, problemsOnly, snatched, jsonOutput, stdOutput bool, version string) (*Propolis, string, error) {
	logthis.Info(ui.YellowBold(ArrowHeader+"Analysing "+path), logthis.NORMAL)

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
	analysis := NewPropolis(path, problemsOnly)
	if jsonOutput || !stdOutput {
		analysis.ToggleStdOutput(false)
	}

	// general checks
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking Path is a music release"), logthis.NORMAL)
	err := release.ParseFiles()

	analysis.ErrorCheck(LevelCritical, "2.3.1", "Release contains FLAC files", "Error parsing files", err, DoNotAppendError)
	if err != nil {
		analysis.ConditionCheck(LevelCritical, "2.2.10.8", ArrowHeader+"At least one FLAC has illegal ID3v2 tags.", ArrowHeader+err.Error(), errors.Is(err, flac.ErrNoFlacHeader))
	}
	analysis.ConditionCheck(LevelInfo, internalRule, "Total size of release folder: "+strconv.FormatFloat(totalSize, 'f', 2, 32)+"Mb.", "", true)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleMusic), logthis.NORMAL)
	analysis = CheckMusicFiles(release, analysis)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleOrganization), logthis.NORMAL)
	analysis = CheckOrganization(release, snatched, analysis)

	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleTags), logthis.NORMAL)
	analysis = CheckTags(release, analysis)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFilenames), logthis.NORMAL)
	analysis = CheckFilenames(release, analysis)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleExtraFiles), logthis.NORMAL)
	analysis = CheckExtraFiles(release, analysis)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFoldername), logthis.NORMAL)
	analysis = CheckFolderName(release, analysis)

	var overviewFile string
	if !disableSpecs {
		logthis.Info(titleHeader+ui.BlueBoldUnderlined("Generating spectrograms"), logthis.NORMAL)
		overviewFile, err = GenerateSpectrograms(release, !jsonOutput)
		if err != nil {
			logthis.Error(err, logthis.NORMAL)
		} else {
			logthis.Info(ui.BlueBold("Spectrograms generated in "+metadataDir+". Check for transcodes (see wiki#408)."), logthis.NORMAL)
		}
	}
	if jsonOutput {
		analysis.ToggleStdOutput(true)
		// TODO take --only-problems into account!
		logthis.Info(analysis.JSONOutput(), logthis.NORMAL)
	} else {
		logthis.Info("\n"+titleHeader+ui.BlueBoldUnderlined("Results\n")+ui.Blue(analysis.Summary()), logthis.NORMAL)
	}
	// saving log to file
	if err != analysis.SaveOuput(metadataDir, version) {
		return analysis, overviewFile, err
	}
	return analysis, overviewFile, nil
}
