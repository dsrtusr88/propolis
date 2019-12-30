package main

import (
	"strconv"

	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/music"
)

var (
	nonFlacExtensions = []string{".mp3", ".aac", ".ogg", ".alac", ".opus", ".ac3"}
)

func CheckMusicFiles(release *music.Release) error {
	isConsistent, bitDepth := release.CheckConsistentBitDepth()
	log.NonCriticalResult(isConsistent, internalRule, "All files are "+bitDepth+"bit files.", "Release has a mix of bit depth.")
	if !isConsistent {
		log.BadResult(release.Has24bitTracks(), internalRule, "At least one track is 24bit FLAC when the rest is 16bit.", "Inconsistent bit depts but no 24bit track.")
		// TODO check inconsistent but > 24bit
	} else {
		bitD, _ := strconv.Atoi(bitDepth)
		log.CriticalResult(bitD <= 24, "2.1.1", "All bit depths are less than 24bit. ", "Bit depths exceeding maximum of 24.")
	}

	isConsistent, sampleRate := release.CheckConsistentSampleRate()
	log.NonCriticalResult(isConsistent, internalRule, "All files have a sample rate of "+sampleRate+"Hz.", "Release has a mix of sample rates.")
	if isConsistent {
		sr, _ := strconv.Atoi(sampleRate)
		log.CriticalResult(sr <= 192000, "2.1.1", "All sample rates are less than or equal to 192kHz.", "Sample rates exceeding maximum of 192kHz.")
	}
	// TODO if !consistent, check highest sample rate

	// NOTE: is the rule track-by-track or on average in the release? what about the stupid "silent" tracks in some releases before a hidden song?
	minAvgBitRate, maxAvgBitRate := release.CheckMinMaxBitrates()
	log.CriticalResult(minAvgBitRate > 192000, "2.1.3", "All tracks have at least 192kbps bitrate (between "+strconv.Itoa(minAvgBitRate)+" and "+strconv.Itoa(maxAvgBitRate)+").", "At least one file has a lower than 192kbps bit rate: "+strconv.Itoa(minAvgBitRate))

	// checking for mutt rip
	for _, ext := range nonFlacExtensions {
		files, err := fs.GetFilesByExt(release.Path, ext)
		if err != nil {
			log.BadResult(err == nil, internalRule, "", "Critical error: "+err.Error())
			return err
		}
		log.CriticalResult(len(files) == 0, "2.1.6.3", "Release does not also contain "+ext+" files.", "Release also contains "+ext+" files, possible mutt rip.")
	}

	// TODO check for ID3 tags 2.2.10.8
	return nil
}
