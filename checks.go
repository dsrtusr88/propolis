package propolis

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.com/catastrophic/assistance/flac"
	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/music"
	"gitlab.com/catastrophic/assistance/strslice"
)

func CheckMusicFiles(release *music.Release, res *Propolis) *Propolis {
	// running the checks
	err := release.CheckVendor()
	res.ErrorCheck(LevelWarning, "2.1.6", OKSameEncoder, KOSameEncoder, err, AppendError)

	isConsistent, bitDepth := release.CheckConsistentBitDepth()
	res.ConditionCheck(LevelWarning, "2.1.6", fmt.Sprintf(OKSameBitDepth, bitDepth), KOSameBitDepth, isConsistent)
	if !isConsistent {
		res.ConditionCheck(LevelAwful, "2.1.6.2", ArrowHeader+OKOne24bitTrack, ArrowHeader+KOOne24bitTrack, release.Has24bitTracks())
		// TODO check inconsistent but > 24bit
	} else {
		bitD, _ := strconv.Atoi(bitDepth)
		res.ConditionCheck(LevelCritical, "2.1.1", ArrowHeader+OKValidBitDepth, ArrowHeader+KOValidBitDepth, bitD <= 24)
	}

	isConsistent, sampleRate := release.CheckConsistentSampleRate()
	res.ConditionCheck(LevelWarning, "2.1.6", fmt.Sprintf(OKSameSampleRate, sampleRate), KOSameSampleRate, isConsistent)
	sr, _ := strconv.Atoi(sampleRate)
	res.ConditionCheck(LevelCritical, "2.1.1", ArrowHeader+OKValidSampleRate, ArrowHeader+KOValidSampleRate, sr <= 192000)

	// NOTE: is the rule track-by-track or on average in the release? what about the stupid "silent" tracks in some releases before a hidden song?
	minAvgBitRate, maxAvgBitRate := release.CheckMinMaxBitrates()
	res.ConditionCheck(LevelCritical, "2.1.3", "All tracks have at least 192kbps bitrate (between "+strconv.Itoa(minAvgBitRate/1000)+"kbps and "+strconv.Itoa(maxAvgBitRate/1000)+"kbps).", "At least one file has a lower than 192kbps bit rate: "+strconv.Itoa(minAvgBitRate), minAvgBitRate > 192000)

	// checking for mutt rip
	forbidden := fs.GetAllowedFilesByExt(release.Path, nonFlacMusicExtensions)
	res.ConditionCheck(LevelCritical, "2.1.6.3", OKMuttRip, fmt.Sprintf(KOMuttRip, strings.Join(forbidden, ",")), len(forbidden) == 0)

	// checking flacs
	err = release.Check()
	res.ErrorCheck(LevelCritical, "2.2.10.8", integrityCheckOK, KOIntegrityCheck, err, DoNotAppendError)
	if err != nil {
		if errors.Is(err, flac.ErrNoFlacHeader) {
			res.ErrorCheck(LevelCritical, "2.2.10.8", "", ArrowHeader+KOID3Tags, err, AppendError)
		} else {
			res.ErrorCheck(LevelCritical, internalRule, "", ArrowHeader+KOIntegrity, err, AppendError)
		}
	}

	// checking for id3v1 tags
	err = release.CheckForID3v1Tags()
	res.ErrorCheck(LevelWarning, internalRule, "No ID3v1 tags detected in the first track.", "The first track contains ID3v1 tags at the end of the file.", err, AppendError)

	// checking for uncompressed flacs
	err = release.CheckCompression()
	res.ErrorCheck(LevelCritical, "2.2.10.10", "First track does not seem to be uncompressed FLAC.", "Error checking for uncompressed FLAC.", err, DoNotAppendError)
	if err != nil {
		if errors.Is(err, flac.ErrorUncompressed) {
			res.ErrorCheck(LevelCritical, "2.2.10.10", "", ArrowHeader+"The first track is uncompressed FLAC", err, AppendError)
		} else {
			res.ErrorCheck(LevelCritical, "2.2.10.10", "", ArrowHeader+"Other error", err, AppendError)
		}
	}
	return res
}

func CheckOrganization(release *music.Release, snatched bool, res *Propolis) *Propolis {
	notTooLong := fs.GetMaxPathLength(release.Path) < 180
	res.ConditionCheck(LevelCritical, "2.3.12", "Maximum character length is less than 180 characters.", "Maximum character length exceeds 180 characters.", notTooLong)
	if !notTooLong {
		for _, f := range fs.GetExceedinglyLongPaths(release.Path, 180) {
			res.ConditionCheck(LevelCritical, "2.3.12", "", ArrowHeader+"Too long: "+f, false)
		}
	}

	// checking for only allowed extensions are used
	if snatched {
		allowedExtensions = append(allowedExtensions, ".json")
	}
	forbidden := fs.GetForbiddenFilesByExt(release.Path, allowedExtensions)
	res.ConditionCheck(LevelCritical, "wiki#371", "Release only contains allowed extensions. ", "Release contains forbidden extensions, which would be rejected by upload.php.", len(forbidden) == 0)
	if len(forbidden) != 0 {
		res.ConditionCheck(LevelCritical, "wiki#371", "", ArrowHeader+"Forbidden files: "+strings.Join(forbidden, ", "), false)
	}

	// checking for empty dirs or uselessly nested folders
	res.ConditionCheck(LevelCritical, "2.3.3", "Release does not have empty folders or unnecessary nested folders.", "Release has empty folders or unnecessary nested folders.", !fs.HasEmptyNestedFolders(release.Path))
	res.ConditionCheck(LevelCritical, "2.3.20", "No leading space/dot found in files and folders.", "Release has files or folders with a leading space or dot.", len(fs.GetFilesAndFoldersByPrefix(release.Path, forbiddenLeadingCharacters)) == 0)
	err := release.CheckMultiDiscOrganization()
	res.ErrorCheck(LevelCritical, "2.3.15", "Release is not multi-disc, or files from multiple discs are either in top folder with disc numbers in filenames, or in dedicated subfolders.", "Tracks from this multi-disc release are incorrectly organized", err, AppendError)

	return res
}

func CheckTags(release *music.Release, res *Propolis) *Propolis {
	res.ErrorCheck(LevelCritical, "2.3.16.1/4", OKRequiredTags, KORequiredTags, release.CheckTags(), DoNotAppendError)
	res.ErrorCheck(LevelCritical, internalRule, OKMetadataSize, KOMetadataSize, release.CheckMaxMetadataSize(Size1024KiB), DoNotAppendError)
	res.ConditionCheck(LevelCritical, "2.3.19", OKCoverSize, KOCoverSize, release.CheckMaxCoverAndPaddingSize() <= Size1024KiB)
	res.ErrorCheck(LevelCritical, internalRule, OKConsistentTags, KOConsistentTags, release.CheckConsistentTags(), AppendError)
	// TODO album title can be different in case of multidisc -- 2.3.18.3.3
	res.ErrorCheck(LevelWarning, internalRule, OKConsistentAlbumArtist, KOConsistentAlbumArtist, release.CheckAlbumArtist(), AppendError)
	// TODO check combined tags
	// TODO export tags to txt file
	return res
}
func CheckFilenames(release *music.Release, res *Propolis) *Propolis {
	// checking for forbidden characters
	withForbiddenChars := fs.GetFilesAndFoldersBySubstring(release.Path, forbiddenCharacters)
	res.ConditionCheck(LevelCritical, internalRule, OKValidCharacters, KOValidCharacters, len(withForbiddenChars) == 0)
	if len(withForbiddenChars) != 0 {
		res.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, ArrowHeader+fmt.Sprintf(InvalidCharacters, strings.Join(withForbiddenChars, ", ")), len(withForbiddenChars) == 0)
	}
	// detecting track.FLAC, track.Flac
	var capitalizedExt bool
	for _, f := range release.Flacs {
		if strings.ToLower(filepath.Ext(f.Path)) == ".flac" && filepath.Ext(f.Path) != ".flac" {
			capitalizedExt = true
			break
		}
	}
	res.ConditionCheck(LevelWarning, internalRule, OKLowerCaseExtensions, KOLowerCaseExtensions, !capitalizedExt)
	// checking filenames contain track numbers and (at least part of) the title
	if len(release.Flacs) != 1 {
		res.ConditionCheck(LevelCritical, "2.3.13", OKTrackNumbersInFilenames, KOTrackNumbersInFilenames, release.CheckTrackNumbersInFilenames())
	} else {
		res.ConditionCheck(LevelWarning, "2.3.13", OKTrackNumberInFilename, KOTrackNumberInFilename, release.CheckTrackNumbersInFilenames())
	}
	res.ConditionCheck(LevelCritical, "2.3.11", OKTitleInFilenames, KOTitleInFilenames, release.CheckFilenameContainsStartOfTitle(minTitleSize))
	// checking filename order
	ordered, err := release.CheckFilenameOrder()
	if err != nil {
		res.ErrorCheck(LevelCritical, internalRule, BlankBecauseImpossible, KOCheckingFilenameOrder, err, AppendError)
	} else {
		res.ConditionCheck(LevelCritical, "2.3.14./.2", OKFilenameOrder, KOFilenameOrder, ordered)
	}
	return res
}

func CheckFolderName(release *music.Release, res *Propolis) *Propolis {
	if len(release.Flacs) == 0 {
		res.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, KOFlacPresent, len(release.Flacs) != 0)
		return res
	}
	// comparisons are case insensitive
	folderName := strings.ToLower(filepath.Base(release.Path))
	// getting metadata
	tags := release.Flacs[0].CommonTags()
	title := tags.Album

	// checking title is in folder name
	res.ConditionCheck(LevelCritical, "2.3.2", OKTitleInFoldername, KOTitleInFoldername, strings.Contains(folderName, strings.ToLower(title)))
	// checking artists are in the folder name
	artists := tags.AlbumArtist
	if len(artists) == 0 {
		// no album artist found, falling back to artists
		for _, f := range release.Flacs {
			artists = append(artists, f.CommonTags().Artist...)
		}
		strslice.RemoveDuplicates(&artists)
	}
	// if more than 3 artists, release should be VA
	if len(artists) >= 3 {
		artists = []string{"Various Artists", "VA"}
	}
	var artistsNotFound []string
	for _, a := range artists {
		if !strings.Contains(folderName, strings.ToLower(a)) {
			artistsNotFound = append(artistsNotFound, a)
		}
	}
	if len(artists) >= 3 && len(artistsNotFound) == 1 {
		// one of the "Various Artists" forms was found, considering everything was found
		artistsNotFound = []string{}
	}
	res.ConditionCheck(LevelWarning, "2.3.2", OKArtistsInFoldername, fmt.Sprintf(KOArtistsInFoldername, strings.Join(artistsNotFound, ", ")), len(artistsNotFound) == 0)
	// checking year is mentioned
	year := tags.Year
	date := tags.Date
	if year != "" || date != "" {
		var foundYear, foundDate bool
		if year != "" {
			foundYear = strings.Contains(folderName, year)
		}
		if date != "" {
			foundDate = strings.Contains(folderName, date)
		}
		res.ConditionCheck(LevelWarning, "2.3.2", OKYearInFoldername, KOYearInFoldername, foundYear || foundDate)
	}
	// checking if formal is mentioned
	res.ConditionCheck(LevelWarning, "2.3.2", OKFormatInFoldername, KOFormatInFoldername, strings.Contains(folderName, "flac"))
	if release.Has24bitTracks() {
		res.ConditionCheck(LevelWarning, "2.3.2", OK24BitInFoldername, KO24BitInFoldername, strings.Contains(folderName, "24"))
	}
	// checking if source is mentioned
	logsAndCues := fs.GetAllowedFilesByExt(release.Path, []string{".log", ".cue"})
	if len(logsAndCues) != 0 {
		res.ConditionCheck(LevelWarning, "2.3.2", OKCDInFoldername, KOCDInFoldername, strings.Contains(folderName, "cd"))
	} else {
		res.ConditionCheck(LevelWarning, "2.3.2", OKWEBInFoldername, KOWEBInFoldername, strings.Contains(folderName, "web"))
	}

	return res
}

func CheckExtraFiles(release *music.Release, res *Propolis) *Propolis {
	// checking for cover
	res.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKCoverFound, music.DefaultCover), fmt.Sprintf(KOCoverFound, music.DefaultCover), release.HasCover())
	// checking for extra files
	nonMusic := fs.GetAllowedFilesByExt(release.Path, nonMusicExtensions)
	res.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKExtraFiles, len(nonMusic)), KOExtraFiles, len(nonMusic) != 0)
	// displaying extra files size and checking ratio vs. music files
	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)
	nonMusicSize := float64(fs.GetPartialSize(release.Path, nonMusic)) / (1024 * 1024)
	ratio := 100 * nonMusicSize / totalSize
	res.ConditionCheck(LevelInfo, internalRule, fmt.Sprintf(OKExtraFilesSize, strconv.FormatFloat(nonMusicSize, 'f', 2, 32)), BlankBecauseImpossible, true)
	res.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKExtraFilesRatio, strconv.FormatFloat(ratio, 'f', 2, 32)), fmt.Sprintf(KOExtraFilesRatio, strconv.FormatFloat(ratio, 'f', 2, 32)), ratio < 10)
	return res
}

func GenerateSpectrograms(release *music.Release, verbose bool) (string, error) {
	// generating full spectrograms
	_, err := release.GenerateSpectrograms("propolis", verbose)
	if err != nil {
		return "", err
	}
	// combination of 10s slices from each song
	return release.GenerateCombinedSpectrogram(verbose)
}
