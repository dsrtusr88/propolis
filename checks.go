package propolis

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"gitlab.com/catastrophic/assistance/flac"
	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/music"
	"gitlab.com/catastrophic/assistance/strslice"
)

func (p *Propolis) CheckRelease() {
	totalSize := float64(fs.GetTotalSize(p.release.Path)) / (1024 * 1024)
	err := p.release.ParseFiles()
	p.ErrorCheck(LevelCritical, "2.3.1", OKReleaseHasFlacs, KOReleaseHasFlacs, err, DoNotAppendError)
	if err != nil {
		p.ConditionCheck(LevelTrulyAwful, "2.2.10.8", ArrowHeader+KOID3v2Tags, ArrowHeader+err.Error(), errors.Is(err, flac.ErrorID3v2Header))
	}
	p.ConditionCheck(LevelInfo, internalRule, fmt.Sprintf(OKTotalSize, strconv.FormatFloat(totalSize, 'f', 2, 32)), BlankBecauseImpossible, true)
	if len(p.release.Flacs) == 0 {
		p.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, KONoTracks, len(p.release.Flacs) != 0)
	}
}

func (p *Propolis) CheckMusicFiles() {
	// checking the encoder
	p.ErrorCheck(LevelWarning, "2.1.6", OKSameEncoder, KOSameEncoder, p.release.CheckVendor(), AppendError)
	// checking for consistency in bit depth
	isConsistent, bitDepth := p.release.CheckConsistentBitDepth()
	p.ConditionCheck(LevelWarning, "2.1.6", fmt.Sprintf(OKSameBitDepth, bitDepth), KOSameBitDepth, isConsistent)
	if !isConsistent {
		p.ConditionCheck(LevelAwful, "2.1.6.2", ArrowHeader+OKOne24bitTrack, ArrowHeader+KOOne24bitTrack, p.release.Has24bitTracks())
		// TODO check inconsistent but > 24bit
	} else {
		bitD, _ := strconv.Atoi(bitDepth)
		p.ConditionCheck(LevelCritical, "2.1.1", ArrowHeader+OKValidBitDepth, ArrowHeader+KOValidBitDepth, bitD <= 24)
	}
	// checking for consistency in sample rate
	isConsistent, sampleRate := p.release.CheckConsistentSampleRate()
	p.ConditionCheck(LevelWarning, "2.1.6", fmt.Sprintf(OKSameSampleRate, sampleRate), KOSameSampleRate, isConsistent)
	sr, _ := strconv.Atoi(sampleRate)
	p.ConditionCheck(LevelCritical, "2.1.1", ArrowHeader+OKValidSampleRate, ArrowHeader+KOValidSampleRate, sr <= 192000)
	// NOTE: is the rule track-by-track or on average in the release? what about the stupid "silent" tracks in some releases before a hidden song?
	minAvgBitRate, maxAvgBitRate := p.release.CheckMinMaxBitrates()
	p.ConditionCheck(LevelCritical, "2.1.3", fmt.Sprintf(OKBitRate, strconv.Itoa(minAvgBitRate/1000), strconv.Itoa(maxAvgBitRate/1000)), fmt.Sprintf(KOBitRate, strconv.Itoa(minAvgBitRate/1000)), minAvgBitRate > 192000)
	// checking if mutt rip
	forbidden := fs.GetAllowedFilesByExt(p.release.Path, nonFlacMusicExtensions)
	p.ConditionCheck(LevelCritical, "2.1.6.3", OKMuttRip, fmt.Sprintf(KOMuttRip, strings.Join(forbidden, ",")), len(forbidden) == 0)
	// checking flac integrity
	err := p.release.Check()
	p.ErrorCheck(LevelCritical, "2.2.10.8", integrityCheckOK, KOIntegrityCheck, err, DoNotAppendError)
	if err != nil {
		if errors.Is(err, flac.ErrNoFlacHeader) {
			p.ErrorCheck(LevelCritical, "2.2.10.8", "", ArrowHeader+KOID3Tags, err, AppendError)
		} else {
			p.ErrorCheck(LevelCritical, internalRule, "", ArrowHeader+KOIntegrity, err, AppendError)
		}
	}
	// checking for id3v1 tags
	err = p.release.CheckForID3v1Tags()
	p.ErrorCheck(LevelWarning, internalRule, OKID3v1Tags, KOID3v1Tags, err, DoNotAppendError)
	// checking for uncompressed flacs
	err = p.release.CheckCompression()
	p.ErrorCheck(LevelCritical, "2.2.10.10", OKUncompressedFlac, KOUncompressedFlac, err, DoNotAppendError)
	if err != nil {
		if errors.Is(err, flac.ErrorUncompressed) {
			p.ErrorCheck(LevelCritical, "2.2.10.10", BlankBecauseImpossible, ArrowHeader+KOUncompressed, err, AppendError)
		} else {
			p.ErrorCheck(LevelCritical, "2.2.10.10", BlankBecauseImpossible, ArrowHeader+OtherError, err, AppendError)
		}
	}
	// checking for MQA encoding
	if len(p.release.Flacs) != 0 {
		p.ConditionCheck(LevelCritical, "upload#DNU", OKNoMQAMetadata, KONoMQAMetadata, !p.release.Flacs[0].CheckForMQAMetadata())
		isMQA, _ := p.release.Flacs[0].CheckForMQASyncword()
		p.ConditionCheck(LevelCritical, "upload#DNU", OKNoMQASyncword, KONoMQASyncword, !isMQA)
	}
}

func (p *Propolis) CheckOrganization(snatched bool) {
	// checking for overly long paths
	longFiles := fs.GetExceedinglyLongPaths(p.release.Path, 180)
	if snatched {
		longFiles = IgnoreVarroaFiles(longFiles)
	}
	p.ConditionCheck(LevelCritical, "2.3.12", OKMaxCharacterLength, KOMaxCharacterLength, len(longFiles) == 0)
	if len(longFiles) != 0 {
		for _, f := range longFiles {
			p.ConditionCheck(LevelCritical, "2.3.12", "", ArrowHeader+fmt.Sprintf(KOTooLong, utf8.RuneCountInString(f), f), false)
		}
	}
	// checking for non-standard spaces
	filesWithNonStandardSpaces := fs.GetPathsWithNonStandardSpaces(p.release.Path)
	p.ConditionCheck(LevelWarning, internalRule, OKNonStandardSpaces, KONonStandardSpaces, len(filesWithNonStandardSpaces) == 0)
	if len(filesWithNonStandardSpaces) != 0 {
		for _, f := range filesWithNonStandardSpaces {
			p.ConditionCheck(LevelWarning, internalRule, "", ArrowHeader+fmt.Sprintf(KOIrregularSpaces, f), false)
		}
	}
	// checking for only allowed extensions are used
	forbidden := fs.GetForbiddenFilesByExt(p.release.Path, allowedExtensions)
	if snatched {
		forbidden = IgnoreVarroaFiles(forbidden)
	}
	p.ConditionCheck(LevelCritical, "wiki#371", OKAllowedExtensions, KOAllowedExtensions, len(forbidden) == 0)
	if len(forbidden) != 0 {
		p.ConditionCheck(LevelCritical, "wiki#371", "", ArrowHeader+fmt.Sprintf(KOForbiddenFiles, strings.Join(forbidden, ", ")), false)
	}
	// checking for empty dirs or uselessly nested folders
	p.ConditionCheck(LevelCritical, "2.3.3", OKEmptyFolders, KOEmptyFolders, !fs.HasEmptyNestedFolders(p.release.Path))
	p.ConditionCheck(LevelCritical, "2.3.20", OKNoLeadingDot, KONoLeadingDot, len(fs.GetFilesAndFoldersByPrefix(p.release.Path, forbiddenLeadingCharacters)) == 0)
	err := p.release.CheckMultiDiscOrganization()
	p.ErrorCheck(LevelCritical, "2.3.15", OKMultiDiscOrganization, KOMultiDiscOrganization, err, AppendError)
}

func (p *Propolis) CheckTags() {
	if len(p.release.Flacs) == 0 {
		p.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, KOFlacPresent, len(p.release.Flacs) != 0)
		return
	}

	errs := p.release.CheckMinimalTags()
	p.ConditionCheck(LevelCritical, "2.3.16.1/4", OKRequiredTags, KORequiredTags, len(errs) == 0)
	for _, e := range errs {
		p.ErrorCheck(LevelCritical, "2.3.16.1/4", BlankBecauseImpossible, ArrowHeader+"Error", e, AppendError)
	}
	errs = p.release.CheckMaxMetadataSize(Size1024KiB)
	p.ConditionCheck(LevelCritical, internalRule, OKMetadataSize, KOMetadataSize, len(errs) == 0)
	for _, e := range errs {
		p.ErrorCheck(LevelCritical, internalRule, BlankBecauseImpossible, ArrowHeader+"Error", e, AppendError)
	}
	p.ConditionCheck(LevelCritical, "2.3.19", OKCoverSize, KOCoverSize, p.release.CheckMaxCoverAndPaddingSize() <= Size1024KiB)
	p.ErrorCheck(LevelCritical, internalRule, OKConsistentTags, KOConsistentTags, p.release.CheckConsistentTags(), AppendError)
	// TODO album title can be different in case of multidisc -- 2.3.18.3.3
	p.ErrorCheck(LevelWarning, internalRule, OKConsistentAlbumArtist, KOConsistentAlbumArtist, p.release.CheckAlbumArtist(), AppendError)
	// TODO check combined tags
	p.ConditionCheck(LevelWarning, "2.3.18.3", OKCombinedTrackNumber, KOCombinedTrackNumber, p.release.Flacs[0].CheckNotCombinedTrackNumber())
}

func (p *Propolis) CheckFilenames(snatched bool) {
	// checking for forbidden characters
	withForbiddenChars := fs.GetFilesAndFoldersBySubstring(p.release.Path, forbiddenCharacters)
	if snatched {
		withForbiddenChars = IgnoreVarroaFiles(withForbiddenChars)
	}
	p.ConditionCheck(LevelCritical, internalRule, OKValidCharacters, KOValidCharacters, len(withForbiddenChars) == 0)
	if len(withForbiddenChars) != 0 {
		p.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, ArrowHeader+fmt.Sprintf(InvalidCharacters, strings.Join(withForbiddenChars, ", ")), len(withForbiddenChars) == 0)
	}
	// detecting track.FLAC, track.Flac
	var capitalizedExt bool
	for _, f := range p.release.Flacs {
		if strings.ToLower(filepath.Ext(f.Path)) == ".flac" && filepath.Ext(f.Path) != ".flac" {
			capitalizedExt = true
			break
		}
	}
	p.ConditionCheck(LevelWarning, internalRule, OKLowerCaseExtensions, KOLowerCaseExtensions, !capitalizedExt)
	// checking filenames contain track numbers and (at least part of) the title
	if len(p.release.Flacs) != 1 {
		p.ConditionCheck(LevelCritical, "2.3.13", OKTrackNumbersInFilenames, KOTrackNumbersInFilenames, p.release.CheckTrackNumbersInFilenames())
	} else {
		p.ConditionCheck(LevelWarning, "2.3.13", OKTrackNumberInFilename, KOTrackNumberInFilename, p.release.CheckTrackNumbersInFilenames())
	}
	p.ConditionCheck(LevelCritical, "2.3.11", OKTitleInFilenames, KOTitleInFilenames, p.release.CheckFilenameContainsStartOfTitle(minTitleSize))
	// checking filename order with disc info
	ordered, err := p.release.CheckFilenameOrder(true)
	if err != nil {
		p.ErrorCheck(LevelCritical, internalRule, BlankBecauseImpossible, KOCheckingFilenameOrder, err, AppendError)
	} else {
		p.ConditionCheck(LevelCritical, "2.3.14./.2", OKFilenameOrder, KOFilenameOrder, ordered)
	}
}

func (p *Propolis) CheckFolderName() {
	if len(p.release.Flacs) == 0 {
		p.ConditionCheck(LevelCritical, internalRule, BlankBecauseImpossible, KOFlacPresent, len(p.release.Flacs) != 0)
		return
	}
	// comparisons are case insensitive
	folderName := strings.ToLower(filepath.Base(p.release.Path))
	// getting metadata
	tags := p.release.Flacs[0].CommonTags()
	title := tags.Album

	// checking title is in folder name
	p.ConditionCheck(LevelCritical, "2.3.2", OKTitleInFoldername, KOTitleInFoldername, flac.StringContainsStartOfAnother(folderName, title, 30))
	// checking artists are in the folder name
	artists := tags.AlbumArtist
	if len(artists) == 0 {
		// no album artist found, falling back to artists
		for _, f := range p.release.Flacs {
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
	p.ConditionCheck(LevelWarning, "2.3.2", OKArtistsInFoldername, fmt.Sprintf(KOArtistsInFoldername, strings.Join(artistsNotFound, ", ")), len(artistsNotFound) == 0)
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
		p.ConditionCheck(LevelWarning, "2.3.2", OKYearInFoldername, KOYearInFoldername, foundYear || foundDate)
	}
	// checking if formal is mentioned
	p.ConditionCheck(LevelWarning, "2.3.2", OKFormatInFoldername, KOFormatInFoldername, strings.Contains(folderName, "flac"))
	if p.release.Has24bitTracks() {
		p.ConditionCheck(LevelWarning, "2.3.2", OK24BitInFoldername, KO24BitInFoldername, strings.Contains(folderName, "24"))
	}
	// checking if source is mentioned
	logsAndCues := fs.GetAllowedFilesByExt(p.release.Path, []string{".log", ".cue"})
	var cleanLogsAndCues []string
	for _, f := range logsAndCues {
		if !strings.Contains(f, "propolis") {
			cleanLogsAndCues = append(cleanLogsAndCues, f)
		}
	}
	if len(cleanLogsAndCues) != 0 {
		p.ConditionCheck(LevelWarning, "2.3.2", OKCDInFoldername, KOCDInFoldername, strings.Contains(folderName, "cd"))
	} else {
		p.ConditionCheck(LevelWarning, "2.3.2", OKWEBInFoldername, KOWEBInFoldername, strings.Contains(folderName, "web") || strings.Contains(folderName, "vinyl"))
	}
}

func (p *Propolis) CheckExtraFiles() {
	// checking for cover
	p.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKCoverFound, music.DefaultCover), fmt.Sprintf(KOCoverFound, music.DefaultCover), p.release.HasCover())
	// checking for extra files
	nonMusic := fs.GetAllowedFilesByExt(p.release.Path, nonMusicExtensions)
	p.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKExtraFiles, len(nonMusic)), KOExtraFiles, len(nonMusic) != 0)
	// displaying extra files size and checking ratio vs. music files
	totalSize := float64(fs.GetTotalSize(p.release.Path)) / (1024 * 1024)
	nonMusicSize := float64(fs.GetPartialSize(p.release.Path, nonMusic)) / (1024 * 1024)
	ratio := 100 * nonMusicSize / totalSize
	p.ConditionCheck(LevelInfo, internalRule, fmt.Sprintf(OKExtraFilesSize, strconv.FormatFloat(nonMusicSize, 'f', 2, 32)), BlankBecauseImpossible, true)
	p.ConditionCheck(LevelWarning, internalRule, fmt.Sprintf(OKExtraFilesRatio, strconv.FormatFloat(ratio, 'f', 2, 32)), fmt.Sprintf(KOExtraFilesRatio, strconv.FormatFloat(ratio, 'f', 2, 32)), ratio < 10)
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
