package propolis

import (
	"path/filepath"

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

	// creating overall check struct and adding the first checks
	analysis := NewPropolis(path, release, problemsOnly)
	if jsonOutput || !stdOutput {
		analysis.ToggleStdOutput(false)
	}

	var overviewFile string
	var err error

	// general checks
	logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleRelease), logthis.NORMAL)
	analysis.CheckRelease()
	if len(analysis.release.Flacs) != 0 {
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleMusic), logthis.NORMAL)
		analysis.CheckMusicFiles()
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleOrganization), logthis.NORMAL)
		analysis.CheckOrganization(snatched)
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleTags), logthis.NORMAL)
		analysis.CheckTags()
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFilenames), logthis.NORMAL)
		analysis.CheckFilenames(snatched)
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleExtraFiles), logthis.NORMAL)
		analysis.CheckExtraFiles()
		logthis.Info(titleHeader+ui.BlueBoldUnderlined(TitleFoldername), logthis.NORMAL)
		analysis.CheckFolderName()

		if !disableSpecs {
			logthis.Info(titleHeader+ui.BlueBoldUnderlined("Generating spectrograms"), logthis.NORMAL)
			overviewFile, err = GenerateSpectrograms(release, !stdOutput || !jsonOutput)
			if err != nil {
				logthis.Error(err, logthis.NORMAL)
			} else {
				logthis.Info(ui.BlueBold("Spectrograms generated in "+metadataDir+". Check for transcodes (see wiki#408)."), logthis.NORMAL)
			}
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
