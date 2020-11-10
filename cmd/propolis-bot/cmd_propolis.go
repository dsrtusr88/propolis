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
	"gitlab.com/catastrophic/assistance/strslice"
	"gitlab.com/passelecasque/obstruction/tracker"
)

func sendToVarroa(command *bot.Cmd) (string, error) {
	// TODO extract torrentID, for now assuming it's just the tID
	torrentID, err := strconv.Atoi(command.RawArgs)
	if err != nil {
		return "", errors.New("invalid torrent ID")
	}

	// TODO MAKE A GLOBAL TRACKER
	b, err := tracker.NewGazelle(config.Varroa.Site, config.Varroa.URL, "", "", "", "", config.Varroa.APIKey, userAgent())
	if err != nil {
		return "", err
	}
	b.StartRateLimiter()
	if err := b.Login(); err != nil {
		return "", err
	}
	// getting torrent information
	gTorrent, err := b.GetTorrent(torrentID)
	if err != nil {
		return "", errors.Wrap(err, "invalid torrent ID")
	}
	// checking format
	if gTorrent.Torrent.Format != tracker.FormatFLAC {
		return "", errors.New("this torrent is not FLAC")
	}

	// checking the content looks tolerable
	if strslice.Contains(config.Varroa.BlacklistedUploaders, gTorrent.Torrent.Username) {
		return "", errors.New("torrent uploader is blacklisted")
	}
	for _, tag := range gTorrent.Group.Tags {
		if strslice.Contains(config.Varroa.ExcludedTags, tag) {
			return "", errors.New("torrent content is blacklisted")
		}
	}
	// printing a short description
	var artists []string
	for _, a := range gTorrent.Group.MusicInfo.Artists {
		artists = append(artists, a.Name)
	}

	// sending varroa an order
	q := url.Values{}
	q.Set("site", config.Varroa.Site)
	q.Set("token", config.Varroa.Token)
	q.Set("id", strconv.Itoa(torrentID))
	reqURL := fmt.Sprintf("http://localhost:%d/get/%d?%s", config.Varroa.Port, torrentID, q.Encode())
	logthis.Info(reqURL, logthis.NORMAL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	// TODO have a global httpclient
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
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
