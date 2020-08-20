package propolis

const (
	minTitleSize = 10
	Size1024KiB  = 1024 * 1024

	TitleMusic        = "Checking music files"
	TitleOrganization = "Checking organization"
	TitleTags         = "Checking tags"
	TitleFilenames    = "Checking filenames"
	TitleExtraFiles   = "Checking extra files"
	TitleFoldername   = "Checking folder name"

	BlankBecauseImpossible = ""

	OKSameEncoder     = "The same encoder was used for all tracks."
	KOSameEncoder     = "Could not confirm the same encoder was used"
	OKSameBitDepth    = "All files are %sbit files."
	KOSameBitDepth    = "The tracks do not have the same bit depth."
	OKOne24bitTrack   = "At least one track is 24bit FLAC when the rest is 16bit, acceptable for some WEB releases."
	KOOne24bitTrack   = "Inconsistent bit depths but no 24bit track."
	OKValidBitDepth   = "All bit depths are less than 24bit."
	KOValidBitDepth   = "Bit depths exceeding maximum of 24."
	OKSameSampleRate  = "All files have a sample rate of %sHz."
	KOSameSampleRate  = "Release has a mix of sample rates, acceptable for some WEB releases (2.1.6.2)."
	OKValidSampleRate = "All sample rates are less than or equal to 192kHz."
	KOValidSampleRate = "Sample rates exceeding maximum of 192kHz."

	OKMuttRip        = "Release does not also contain other kinds of music files."
	KOMuttRip        = "Release also contains other music formats, possible mutt rip: %s"
	KOIntegrityCheck = "At least one track is not a valid FLAC file."
	KOID3Tags        = "At least one FLAC has illegal ID3 tags"
	KOIntegrity      = "At least one FLAC has failed an integrity test"

	OKRequiredTags          = "All tracks have at least the required tags."
	KORequiredTags          = "At least one tracks is missing required tags."
	OKMetadataSize          = "All tracks have metadata blocks of a total size smaller then 1024 KiB."
	KOMetadataSize          = "At least one track has metadata blocks of a total size bigger than 1024 KiB, probably due to excessive padding, embedded art, or more exotic things."
	OKCoverSize             = "All tracks either have no embedded art, or the embedded art size is less than 1024KiB (padding included)."
	KOCoverSize             = "At least one track has embedded art and padding exceeding the maximum allowed size of 1024 KiB."
	OKConsistentTags        = "Release-level tags seem consistent among tracks."
	KOConsistentTags        = "Tracks have inconsistent tags about the release."
	OKConsistentAlbumArtist = "Artist/Album artist tags seem consistent."
	KOConsistentAlbumArtist = "Artist/Album artist tags differ from file to file"

	OKValidCharacters         = "Tracks filenames do not appear to contain problematic characters."
	KOValidCharacters         = "At least one track filename or folder contains problematic characters."
	InvalidCharacters         = "In files and folders: %s"
	OKLowerCaseExtensions     = "Track filenames have lower case extensions."
	KOLowerCaseExtensions     = "At least one filename has an uppercase .FLAC extension."
	OKTrackNumbersInFilenames = "All tracks filenames appear to contain their track number."
	KOTrackNumbersInFilenames = "At least one track filename does not contain its track number."
	OKTrackNumberInFilename   = "The track filename appears to contain the track number."
	KOTrackNumberInFilename   = "The track filename does not contain the track number. It is not required for singles, but good practice nonetheless."
	OKTitleInFilenames        = "All tracks filenames appear to contain at least the beginning of song titles."
	KOTitleInFilenames        = "At least one track filename does not seem to include the beginning of the song title."
	OKFilenameOrder           = "Files and subfolder names respect the playing order of the release."
	KOFilenameOrder           = "Files and/or subfolder names do not sort alphabetically into the playing order of the release."
	KOCheckingFilenameOrder   = "Could not check filename/subfolder order. Track/Disc numbers might not be numbers"
	KOFlacPresent             = "Release has no FLACs!"
	OKTitleInFoldername       = "Title of album is in folder name."
	KOTitleInFoldername       = "Title of album (as found in the tags of the first track) is not in the folder name."
	OKArtistsInFoldername     = "All album artists found in folder name."
	KOArtistsInFoldername     = "Not all album artists (as found in the tags of the first track) found in the folder name: missing %s"
	OKYearInFoldername        = "Year of album is in folder name."
	KOYearInFoldername        = "Year of album (as found in the tags of the first track) is not in the folder name."
	OKFormatInFoldername      = "Format (FLAC) found in folder name."
	KOFormatInFoldername      = "Format (FLAC) not found in folder name."
	OK24BitInFoldername       = "Folder name properly mentions the release has 24bit FLAC tracks."
	KO24BitInFoldername       = "Since release seems to contain 24bit FLACs, the folder name could mention it."
	OKCDInFoldername          = "Release contains .log/.cue files and the folder name properly mentions a CD source."
	KOCDInFoldername          = "Since release contains .log/.cue, it seems to be sourced from CD. The folder name could mention it."
	OKWEBInFoldername         = "Release does not contain .log/.cue files and the folder name properly mentions a WEB source."
	KOWEBInFoldername         = "Since release does not .log/.cue, it is probably sources from WEB. The folder name could mention it."
	OKCoverFound              = "Release has a conventional %s in the top folder or in all disc subfolders."
	KOCoverFound              = "Cannot find %s in top folder or in all disc subfolders, consider adding one or renaming the cover to that name."
	OKExtraFiles              = "Release has %d accompanying files."
	KOExtraFiles              = "Release does not have any kind of accompanying files. Suggestion: consider adding at least a cover."
	OKExtraFilesSize          = "Total size of accompanying files: %sMb."
	OKExtraFilesRatio         = "Accompanying files represent %s%% of the total size."
	KOExtraFilesRatio         = "Accompanying files represent %s%% of the total size. Suggestion: if this is because of high resolution artwork or notes, consider uploading separately and linking the files in the description."
)

var (
	// https://redacted.ch/wiki.php?action=article&id=371
	allowedExtensions      = []string{".ac3", ".accurip", ".azw3", ".chm", ".cue", ".djv", ".djvu", ".doc", ".dmg", ".dts", ".epub", ".ffp", ".flac", ".gif", ".htm", ".html", ".jpeg", ".jpg", ".lit", ".log", ".m3u", ".m3u8", ".m4a", ".m4b", ".md5", ".mobi", ".mp3", ".mp4", ".nfo", ".pdf", ".pls", ".png", ".rtf", ".sfv", ".txt"}
	nonFlacMusicExtensions = []string{".ac3", ".dts", ".m4a", ".m4b", ".mp3", ".mp4", ".aac", ".alac", ".ogg", ".opus"}
	nonMusicExtensions     = []string{".accurip", ".azw3", ".chm", ".cue", ".djv", ".djvu", ".doc", ".dmg", ".epub", ".ffp", ".gif", ".htm", ".html", ".jpeg", ".jpg", ".lit", ".log", ".m3u", ".m3u8", ".md5", ".mobi", ".nfo", ".pdf", ".pls", ".png", ".rtf", ".sfv", ".txt"}

	forbiddenCharacters        = []string{":", "*", `\`, "?", `"`, `<`, `>`, "|", "$", "`"}
	forbiddenLeadingCharacters = []string{" ", "."}
)