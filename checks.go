package main

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/music"
)

const (
	minTitleSize = 10
)

var (
	// https://redacted.ch/wiki.php?action=article&id=371
	allowedExtensions      = []string{".ac3", ".accurip", ".azw3", ".chm", ".cue", ".djv", ".djvu", ".doc", ".dmg", ".dts", ".epub", ".ffp", ".flac", ".gif", ".htm", ".html", ".jpeg", ".jpg", ".lit", ".log", ".m3u", ".m3u8", ".m4a", ".m4b", ".md5", ".mobi", ".mp3", ".mp4", ".nfo", ".pdf", ".pls", ".png", ".rtf", ".sfv", ".txt"}
	nonFlacMusicExtensions = []string{".ac3", ".dts", ".m4a", ".m4b", ".mp3", ".mp4", ".aac", ".alac", ".ogg", ".opus"}
	nonMusicExtensions     = []string{".accurip", ".azw3", ".chm", ".cue", ".djv", ".djvu", ".doc", ".dmg", ".epub", ".ffp", ".gif", ".htm", ".html", ".jpeg", ".jpg", ".lit", ".log", ".m3u", ".m3u8", ".md5", ".mobi", ".nfo", ".pdf", ".pls", ".png", ".rtf", ".sfv", ".txt"}

	forbiddenCharacters        = []string{":", "*", `\`, "?", `"`, `<`, `>`, "|", "$", "`"}
	forbiddenLeadingCharacters = []string{" ", "."}
)

func CheckMusicFiles(release *music.Release) error {
	isConsistent, bitDepth := release.CheckConsistentBitDepth()
	log.NonCriticalResult(isConsistent, "2.1.6", "All files are "+bitDepth+"bit files.", "The tracks do not have the same bit depth.")
	if !isConsistent {
		log.BadResult(release.Has24bitTracks(), "2.1.6.2", "⮕ At least one track is 24bit FLAC when the rest is 16bit, acceptable for some WEB releases.", "⮕ Inconsistent bit depths but no 24bit track.")
		// TODO check inconsistent but > 24bit
	} else {
		bitD, _ := strconv.Atoi(bitDepth)
		log.CriticalResult(bitD <= 24, "2.1.1", "⮕ All bit depths are less than 24bit. ", "⮕ Bit depths exceeding maximum of 24.")
	}

	isConsistent, sampleRate := release.CheckConsistentSampleRate()
	log.NonCriticalResult(isConsistent, "2.1.6", "All files have a sample rate of "+sampleRate+"Hz.", "Release has a mix of sample rates, acceptable for some WEB releases (2.1.6.2).")
	if isConsistent {
		sr, _ := strconv.Atoi(sampleRate)
		log.CriticalResult(sr <= 192000, "2.1.1", "⮕ All sample rates are less than or equal to 192kHz.", "⮕ Sample rates exceeding maximum of 192kHz.")
	}
	// TODO if !consistent, check highest sample rate

	// NOTE: is the rule track-by-track or on average in the release? what about the stupid "silent" tracks in some releases before a hidden song?
	minAvgBitRate, maxAvgBitRate := release.CheckMinMaxBitrates()
	log.CriticalResult(minAvgBitRate > 192000, "2.1.3", "All tracks have at least 192kbps bitrate (between "+strconv.Itoa(minAvgBitRate/1000)+"kbps and "+strconv.Itoa(maxAvgBitRate/1000)+"kbps).", "At least one file has a lower than 192kbps bit rate: "+strconv.Itoa(minAvgBitRate))

	// checking for mutt rip
	forbidden := fs.GetAllowedFilesByExt(release.Path, nonFlacMusicExtensions)
	log.CriticalResult(len(forbidden) == 0, "2.1.6.3", "Release does not also contain other kinds of music files.", "Release also contains other music formats, possible mutt rip: "+strings.Join(forbidden, ","))

	// checking flacs
	err := release.Check()
	log.CriticalResult(err == nil, "2.2.10.8", "Integrity checks for all FLACs OK, no ID3 tags detected.", "At least one track is not a valid FLAC file.")
	if err != nil {
		if err.Error() == music.ErrorContainsID3Tags {
			log.CriticalResult(err == nil, "2.2.10.8", "", "⮕ At least one FLAC has illegal ID3 tags.")
		} else {
			log.CriticalResult(err == nil, internalRule, "", "⮕ At least one FLAC has failed an integrity test.")
		}
	}

	// checking for uncompressed flacs
	err = release.CheckCompression()
	log.CriticalResult(err == nil, "2.2.10.10", "First track does not seem to be uncompressed FLAC.", "Error checking for uncompressed FLAC.")
	if err != nil && err.Error() == music.ErrorUncompressed {
		log.CriticalResult(err == nil, "2.2.10.10", "", "⮕ The first track is uncompressed FLAC.")
	}
	return nil
}

func CheckOrganization(release *music.Release) error {
	notTooLong := fs.GetMaxPathLength(release.Path) < 180
	log.CriticalResult(notTooLong, "2.3.12", "Maximum character length is less than 180 characters.", "Maximum character length exceeds 180 characters.")
	if !notTooLong {
		for _, f := range fs.GetExceedinglyLongPaths(release.Path, 180) {
			log.CriticalResult(notTooLong, "2.3.12", "", "⮕ Too long: "+f)
		}
	}

	// checking for only allowed extensions are used
	forbidden := fs.GetForbiddenFilesByExt(release.Path, allowedExtensions)
	log.CriticalResult(len(forbidden) == 0, "wiki#371", "Release only contains allowed extensions. ", "Release contains forbidden extensions, which would be rejected by upload.php.")
	if len(forbidden) != 0 {
		log.CriticalResult(len(forbidden) == 0, "wiki#371", "", "⮕ Forbidden files: "+strings.Join(forbidden, ", "))
	}

	// checking for empty dirs or uselessly nested folders
	log.CriticalResult(!fs.HasEmptyNestedFolders(release.Path), "2.3.3", "Release does not have empty folders or unnecessary nested folders.", "Release has empty folders or unnecessary nested folders.")

	log.CriticalResult(len(fs.GetFilesAndFoldersByPrefix(release.Path, forbiddenLeadingCharacters)) == 0, "2.3.20", "No leading space/dot found in files and folders.", "Release has files or folders with a leading space or dot.")

	if release.NumberOfDiscs() > 1 {
		err := release.CheckMultiDiscOrganization()
		if err != nil {
			log.CriticalResult(err == nil, "2.3.15", "", "Tracks from this multi-disc release are incorrectly organized: "+err.Error())
		} else {
			log.CriticalResult(err == nil, "2.3.15", "Files from multiple discs are either in top folder with disc numbers in filenames, or in dedicated subfolders.", "")
		}
	}

	return nil
}

func CheckTags(release *music.Release) error {
	log.CriticalResult(release.CheckTags() == nil, "2.3.16.1/4", "All tracks have at least the required tags.", "At least one tracks is missing required tags.")
	log.CriticalResult(release.CheckMaxCoverSize() <= 1024*1024, "2.3.19", "All tracks either have no embedded art, or the embedded art size is less than 1024KiB.", "At least one track has embedded art exceeding the maximum allowed size of 1024 KiB.")

	err := release.CheckConsistentTags()
	log.CriticalResult(err == nil, internalRule, "Release-level tags seem consistent among tracks.", "Tracks have inconsistent tags about the release.")
	if err != nil {
		log.CriticalResult(err == nil, internalRule, "", "⮕ Found "+err.Error())
		// TODO album title can be different in case of multidisc -- 2.3.18.3.3
	}

	// check combined tags

	// TODO export tags to txt file
	return nil
}

func CheckFilenames(release *music.Release) error {
	withForbiddenChars := fs.GetFilesAndFoldersBySubstring(release.Path, forbiddenCharacters)
	log.CriticalResult(len(withForbiddenChars) == 0, internalRule, "Tracks filenames do not appear to contain problematic characters.", "At least one track filename or folder contains problematic characters.")
	if len(withForbiddenChars) != 0 {
		log.CriticalResult(len(withForbiddenChars) == 0, internalRule, "", "⮕ In files and folders: "+strings.Join(withForbiddenChars, ", "))
	}

	var capitalizedExt bool
	for _, f := range release.Flacs {
		if filepath.Ext(f.Filename) == ".FLAC" {
			capitalizedExt = true
			break
		}
	}
	log.NonCriticalResult(!capitalizedExt, internalRule, "Track filenames have lower case extensions.", "At least one filename has an uppercase .FLAC extension.")

	log.CriticalResult(release.CheckTrackNumbersInFilenames(), "2.3.13", "All tracks filenames appear to contain their track number.", "At least one track filename does not contain its track number.")

	log.CriticalResult(release.CheckFilenameContainsStartOfTitle(minTitleSize), "2.3.11", "All tracks filenames appear to contain at least the beginning of song titles.", "At least one track filename does not seem to include the beginning of the song title.")

	ordered, err := release.CheckFilenameOrder()
	if err != nil {
		log.CriticalResult(err == nil, internalRule, "", "Could not check filename/subfolder order. Track/Disc numbers might not be numbers: "+err.Error())
	} else {
		log.CriticalResult(ordered, "2.3.14./.2", "Files and subfolder names respect the playing order of the release.", "Files and/or subfolder names do not sort alphabetically into the playing order of the release.")
	}

	return nil
}

func CheckFolderName(release *music.Release) error {
	if len(release.Flacs) == 0 {
		return errors.New("no tracks")
	}
	// getting metadata
	title := release.Flacs[0].Tags.Album
	log.CriticalResult(strings.Contains(release.Path, title), "2.3.2", "Title of album is in folder name.", "Title of album (as found in the tags of the first track) is not in the folder name.")

	artists := release.Flacs[0].Tags.AlbumArtist
	var artistsFound int
	for _, a := range artists {
		if strings.Contains(release.Path, a) {
			artistsFound++
		}
	}
	log.NonCriticalResult(len(artists) == artistsFound, "2.3.2", "All album artists found in folder name.", "Not all (if any) album artists (as found in the tags of the first track) found in the folder name.")

	year := release.Flacs[0].Tags.Year
	date := release.Flacs[0].Tags.Date
	if year != "" || date != "" {
		var foundYear, foundDate bool
		if year != "" {
			foundYear = strings.Contains(release.Path, year)
		}
		if date != "" {
			foundDate = strings.Contains(release.Path, date)
		}
		log.NonCriticalResult(foundYear || foundDate, "2.3.2", "Year of album is in folder name.", "Year of album (as found in the tags of the first track) is not in the folder name.")
	}
	log.NonCriticalResult(strings.Contains(release.Path, "FLAC"), "2.3.2", "Format (FLAC) found in folder name.", "Format (FLAC) not found in folder name.")

	if release.Has24bitTracks() {
		log.NonCriticalResult(strings.Contains(release.Path, "24"), "2.3.2", "Folder name properly mentions the release has 24bit FLAC tracks.", "Since release seems to contain 24bit FLACs, the folder name could mention it. ")
	}

	logsAndCues := fs.GetAllowedFilesByExt(release.Path, []string{".log", ".cue"})
	if len(logsAndCues) != 0 {
		log.NonCriticalResult(strings.Contains(release.Path, "CD"), "2.3.2", "Release contains .log/.cue files and the folder name properly mentions a CD source.", "Since release contains .log/.cue, it seems to be sourced from CD. The folder name could mention it. ")
	} else {
		log.NonCriticalResult(strings.Contains(release.Path, "WEB"), "2.3.2", "Release does not contain .log/.cue files and the folder name properly mentions a WEB source.", "Since release does not .log/.cue, it is probably sources from WEB. The folder name could mention it. ")
	}

	return nil
}

func CheckExtraFiles(release *music.Release) error {
	log.NonCriticalResult(release.HasCover(), internalRule, "Release has a conventional "+music.DefaultCover+" in the top folder or in all disc subfolders.", "Cannot find "+music.DefaultCover+" in top folder or in all disc subfolders, consider adding one or renaming the cover to that name.")

	nonMusic := fs.GetAllowedFilesByExt(release.Path, nonMusicExtensions)
	log.NonCriticalResult(len(nonMusic) != 0, internalRule, "Release has "+strconv.Itoa(len(nonMusic))+" accompanying files.", "Release does not have any kind of accompanying files. Suggestion: consider adding at least a cover.")

	totalSize := float64(fs.GetTotalSize(release.Path)) / (1024 * 1024)
	nonMusicSize := float64(fs.GetPartialSize(release.Path, nonMusic)) / (1024 * 1024)
	log.NeutralResult(true, internalRule, "Total size of accompanying files: "+strconv.FormatFloat(nonMusicSize, 'f', 2, 32)+"Mb.", "")
	ratio := 100 * nonMusicSize / totalSize
	log.NonCriticalResult(ratio < 10, internalRule, "Accompanying files represent "+strconv.FormatFloat(ratio, 'f', 2, 32)+"% of the total size.", "Accompanying files represent "+strconv.FormatFloat(ratio, 'f', 2, 32)+"% of the total size. Suggestion: if this is because of high resolution artwork or notes, consider uploading separately and linking the files in the description.")

	return nil
}

func GenerateSpectrograms(release *music.Release) ([]string, error) {
	var wg sync.WaitGroup
	var combinedPNG string
	var combinedErr error
	// generating combinedPNG in the background
	wg.Add(1)
	go func() {
		// combination of 10s slices from each song
		combinedPNG, combinedErr = release.GenerateCombinedSpectrogram()
		wg.Done()
	}()
	// generating full spectrograms
	pngs, err := release.GenerateSpectrograms(" ")
	if err != nil {
		return []string{}, err
	}
	// checking combined PNG was correctly created
	wg.Wait()
	if combinedErr != nil {
		return []string{}, combinedErr
	}
	pngs = append([]string{combinedPNG}, pngs...)
	return pngs, nil
}
