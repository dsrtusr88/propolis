package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/music"
)

var (
	// https://redacted.ch/wiki.php?action=article&id=371
	allowedExtensions = []string{".ac3", ".accurip", ".azw3", ".chm", ".cue", ".djv", ".djvu", ".doc", ".dmg", ".dts", ".epub", ".ffp", ".flac", ".gif", ".htm", ".html", ".jpeg", ".jpg", ".lit", ".log", ".m3u", ".m3u8", ".m4a", ".m4b", ".md5", ".mobi", ".mp3", ".mp4", ".nfo", ".pdf", ".pls", ".png", ".rtf", ".sfv", ".txt"}

	nonFlacMusicExtensions = []string{".ac3", ".dts", ".m4a", ".m4b", ".mp3", ".mp4", ".aac", ".alac", ".ogg", ".opus"}

	forbiddenCharacters        = []string{":", "*", `\`, "?", `"`, `<`, `>`, "|", "$", "`"}
	forbiddenLeadingCharacters = []string{" ", "."}
)

func CheckMusicFiles(release *music.Release) error {
	isConsistent, bitDepth := release.CheckConsistentBitDepth()
	log.NonCriticalResult(isConsistent, "2.1.6", "All files are "+bitDepth+"bit files.", "The tracks do not have the same bit depth.")
	if !isConsistent {
		log.BadResult(release.Has24bitTracks(), "2.1.6.2", "At least one track is 24bit FLAC when the rest is 16bit, acceptable for some WEB releases.", "Inconsistent bit depths but no 24bit track.")
		// TODO check inconsistent but > 24bit
	} else {
		bitD, _ := strconv.Atoi(bitDepth)
		log.CriticalResult(bitD <= 24, "2.1.1", "All bit depths are less than 24bit. ", "Bit depths exceeding maximum of 24.")
	}

	isConsistent, sampleRate := release.CheckConsistentSampleRate()
	log.NonCriticalResult(isConsistent, "2.1.6", "All files have a sample rate of "+sampleRate+"Hz.", "Release has a mix of sample rates, acceptable for some WEB releases (2.1.6.2).")
	if isConsistent {
		sr, _ := strconv.Atoi(sampleRate)
		log.CriticalResult(sr <= 192000, "2.1.1", "All sample rates are less than or equal to 192kHz.", "Sample rates exceeding maximum of 192kHz.")
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
			log.CriticalResult(err == nil, "2.2.10.8", "", "At least one FLAC has illegal ID3 tags.")
		} else {
			log.CriticalResult(err == nil, internalRule, "", "At least one FLAC has failed an integrity test.")
		}
	}

	// checking for uncompressed flacs
	err = release.CheckCompression()
	log.CriticalResult(err == nil, "2.2.10.10", "First track does not seem to be uncompressed FLAC.", "Error checking for uncompressed FLAC.")
	if err != nil && err.Error() == music.ErrorUncompressed {
		log.CriticalResult(err == nil, "2.2.10.10", "", "The first track is uncompressed FLAC.")
	}
	return nil
}

func CheckOrganization(release *music.Release) error {
	notTooLong := fs.GetMaxPathLength(release.Path) < 180
	log.CriticalResult(notTooLong, "2.3.12", "Maximum character length is less than 180 characters.", "Maximum character length exceeds 180 characters.")
	if !notTooLong {
		for _, f := range fs.GetExceedinglyLongPaths(release.Path, 180) {
			log.CriticalResult(notTooLong, "2.3.12", "", "Too long: "+f)
		}
	}

	// checking for only allowed extensions are used
	forbidden := fs.GetForbiddenFilesByExt(release.Path, allowedExtensions)
	log.CriticalResult(len(forbidden) == 0, "wiki#371", "Release only contains allowed extensions. ", "Release contains forbidden extensions, which would be rejected by upload.php.")
	if len(forbidden) != 0 {
		log.CriticalResult(len(forbidden) == 0, "wiki#371", "", "Forbidden files: "+strings.Join(forbidden, ", "))
	}

	// checking for empty dirs or uselessly nested folders
	log.CriticalResult(!fs.HasEmptyNestedFolders(release.Path), "2.3.3", "Release does not have empty folders or unnecessary nested folders.", "Release has empty folders or unnecessary nested folders.")

	log.CriticalResult(len(fs.GetFilesAndFoldersByPrefix(release.Path, forbiddenLeadingCharacters)) == 0, "2.3.20", "No leading space/dot found in files and folders.", "Release has files or folders with a leading space or dot.")

	return nil
}

func CheckTags(release *music.Release) error {
	log.CriticalResult(release.CheckTags() == nil, "2.3.16.1/4", "All tracks have at least the required tags.", "At least one tracks is missing required tags.")
	log.CriticalResult(release.CheckMaxCoverSize() <= 1024*1024, "2.3.19", "All tracks either have no embedded art, or the embedded art size is less than 1024KiB.", "At least one track has embedded art exceeding the maximum allowed size of 1024 KiB.")

	err := release.CheckConsistentTags()
	log.CriticalResult(err == nil, internalRule, "Release-level tags seem consistent among tracks.", "Tracks have inconsistent tags about the release.")
	if err != nil {
		log.CriticalResult(err == nil, internalRule, "", "More precisely: "+err.Error())
		// TODO album title can be different in case of multidisc -- 2.3.18.3.3
	}

	// check combined tags

	// TODO export tags to txt file
	return nil
}

func CheckExtraFiles(release *music.Release) error {
	log.NonCriticalResult(fs.FileExists(filepath.Join(release.Path, music.DefaultCover)), internalRule, "Release has a conventional "+music.DefaultCover+" in the top folder.", "Cannot find "+music.DefaultCover+" in top folder, consider adding one or renaming the cover to that name.")

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
