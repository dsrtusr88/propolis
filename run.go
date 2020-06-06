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

func Run(path string, disableSpecs, problemsOnly, snatched bool) (*Results, string, error) {
	// setting output config
	log.problemsOnly = problemsOnly

	// by default, metadata (spectrograms, etc), will be put in a side folder.
	metadataDir := path + " (Metadata)"
	if snatched {
		metadataDir = filepath.Join(path, "Metadata")
	}
	release := music.NewWithExternalMetadata(path, metadataDir)
	res := &Results{}
	var err error

	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking Path is a music release"), logthis.NORMAL)
	err = release.ParseFiles()
	res.Add(log.Critical(err == nil, "2.3.1", "Release contains FLAC files", "Error parsing files"))
	if err != nil {
		if errors.Is(err, flac.ErrNoFlacHeader) {
			res.Add(log.Critical(err == nil, "2.2.10.8", "", arrowHeader+"At least one FLAC has illegal ID3v2 tags."))
		} else {
			res.Add(log.BadResult(err == nil, "2.3.1", "", arrowHeader+"Critical error: "+err.Error()))
		}
		return res, "", err
	}
	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)
	res.Add(log.Info(true, internalRule, "Total size of release folder: "+strconv.FormatFloat(totalSize, 'f', 2, 32)+"Mb.", ""))

	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking music files"), logthis.NORMAL)
	res = CheckMusicFiles(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking organization"), logthis.NORMAL)
	res = CheckOrganization(release, snatched, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking tags"), logthis.NORMAL)
	res = CheckTags(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking filenames"), logthis.NORMAL)
	res = CheckFilenames(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking extra files"), logthis.NORMAL)
	res = CheckExtraFiles(release, res)
	logthis.Info(titleHeader+ui.BlueBoldUnderlined("Checking folder name"), logthis.NORMAL)
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
	logthis.Info("\n"+titleHeader+ui.BlueBoldUnderlined("Results\n")+ui.Blue(res.String()), logthis.NORMAL)

	return res, overviewFile, nil
}
