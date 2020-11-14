package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chat-bot/bot"
	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/flac"
	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/passelecasque/obstruction/tracker"
)

func sendToVarroaDaemon(id int) error {
	q := url.Values{}
	q.Set("site", config.Varroa.Site)
	q.Set("token", config.Varroa.Token)
	q.Set("id", strconv.Itoa(id))
	reqURL := fmt.Sprintf("http://localhost:%d/get/%d?%s", config.Varroa.Port, id, q.Encode())
	logthis.Info(reqURL, logthis.NORMAL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func sendToVarroa(command *bot.Cmd) (string, error) {
	// retrieving the torrent ID
	torrentID, err := tracker.ExtractTorrentID(command.RawArgs)
	if err != nil {
		return "", errors.New("invalid torrent ID")
	}
	// getting torrent information
	gTorrent, err := gazelle.GetTorrent(torrentID)
	if err != nil {
		return "", errors.Wrap(err, "invalid torrent ID")
	}
	// checking format
	if !gTorrent.IsFLAC() {
		return "", errors.New("this torrent is not FLAC")
	}
	// checking the content looks tolerable
	if !gTorrent.IsAcceptable(config.Varroa.BlacklistedUploaders, config.Varroa.ExcludedTags) {
		return "", errors.New("content or uploader is blacklisted")
	}
	// printing a short description
	artists := gTorrent.MainArtists()

	// sending varroa the snatch command
	if err := sendToVarroaDaemon(torrentID); err != nil {
		return "", err
	}

	return fmt.Sprintf("-> %s! Snatching and analyzing torrent \u0002\u000307%s (%d) %s [%s]\u000F for you!", command.User.Nick, flac.SmartArtistList(artists), gTorrent.Torrent.RemasterYear, gTorrent.Group.Name, gTorrent.Torrent.Media), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	bot.RegisterCommand(
		"propolis",
		"Analyze a specific torrent ID",
		"URL or torrentID",
		sendToVarroa)
}
