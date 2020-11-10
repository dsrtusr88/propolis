package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// _ "net/http/pprof"
	"path/filepath"
	"runtime/debug"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/fs"
	"gitlab.com/catastrophic/assistance/imgupload"
	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/privatebin"
	"gitlab.com/passelecasque/obstruction/tracker"
	"gitlab.com/passelecasque/propolis"
)

const (
	WebsocketPort            = 8335
	ErrorNotSnatchWithVarroa = "JSON metadata not found"
)

type IncomingJSON struct {
	Path string
}

func Server() {
	target := config.IRC.CentralBot
	if config.IRC.Role == centralRole {
		target = config.IRC.Channel
	}

	var ptpImgSet bool
	ptpimg, err := imgupload.NewWithAPIKey(config.General.PtpImgKey)
	if err != nil {
		logthis.Error(err, logthis.NORMAL)
	} else {
		ptpImgSet = true
	}

	rtr := mux.NewRouter()
	getMetadata := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		defer debug.FreeOSMemory()
		decoder := json.NewDecoder(r.Body)
		var t IncomingJSON
		err := decoder.Decode(&t)
		if err != nil {
			logthis.Error(err, logthis.NORMAL)
			return
		} else {
			if !fs.DirExists(t.Path) {
				logthis.Info("path not found: "+t.Path, logthis.NORMAL)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if !ircConn.Connected() {
				logthis.Info("not connected to irc, not running propolis.", logthis.NORMAL)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			results, overviewFile, err := propolis.Run(t.Path, false, false, true, false, false, Version)
			if err != nil {
				logthis.Error(err, logthis.NORMAL)
				return
			}

			// sending result to IRC bot
			var link string
			if overviewFile != "" && ptpImgSet {
				link, err = ptpimg.UploadLocalFile(overviewFile)
				if err != nil {
					logthis.Error(err, logthis.NORMAL)
				}
			}

			permalink := ErrorNotSnatchWithVarroa
			var releaseMetadataFile string
			releaseMetadata1 := filepath.Join(t.Path, "TrackerMetadata", "red_release.json")
			releaseMetadata2 := filepath.Join(t.Path, "TrackerMetadata", "red - Release.json")
			if fs.FileExists(releaseMetadata1) {
				releaseMetadataFile = releaseMetadata1
			} else if fs.FileExists(releaseMetadata2) {
				releaseMetadataFile = releaseMetadata2
			}
			if releaseMetadataFile != "" {
				var gt *tracker.GazelleTorrent
				data, err := ioutil.ReadFile(releaseMetadataFile)
				if err != nil {
					logthis.Error(err, logthis.NORMAL)
				} else {
					if err := json.Unmarshal(data, &gt); err != nil {
						logthis.Error(err, logthis.NORMAL)
					} else {
						permalink = "https://redacted.ch/torrents.php?torrentid=" + strconv.Itoa(gt.Torrent.ID)
					}
				}
			}

			results.ToggleStdOutput(true)
			if results.Errors != 0 {
				logthis.Info("problems found for : "+t.Path, logthis.NORMAL)
			}

			var pastedLogURL string
			pastedLogURL, err = privatebin.Post([]byte(results.Output()+"\n-----------\n\n"+results.Tags()), config.General.PrivateBinURL)
			if err != nil {
				logthis.Error(err, logthis.NORMAL)
			}

			// sending summary to IRC chan
			IRCSummary := "\x0309" + filepath.Base(t.Path) + "\x0F | " + permalink + " | " + link + " | " + pastedLogURL + " | " + "\x02\x0307" + results.Summary() + "\x0F"
			ircConn.Privmsg(target, IRCSummary)
			// sending errors one by one, if relevant
			if results.Errors != 0 {
				for _, e := range results.AllErrors() {
					message := []byte("\x0309" + filepath.Base(t.Path) + "\x0F | \x02\x0304Error: " + e + "\x0F")
					ircConn.Privmsg(target, string(message))
				}
			}
		}
	}
	rtr.HandleFunc("/downloads", getMetadata).Methods("POST")
	// rtr.PathPrefix("/debug/").Handler(http.DefaultServeMux)

	// serve
	go func() {
		logthis.Info("Setting up websocket server.", logthis.NORMAL)
		httpServer := &http.Server{Addr: fmt.Sprintf(":%d", WebsocketPort), Handler: rtr}
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logthis.Info("Closing websocket server...", logthis.NORMAL)
			} else {
				logthis.Error(errors.Wrap(err, "error with websocket server"), logthis.NORMAL)
			}
		}
	}()
	logthis.Info("Server is listening.", logthis.NORMAL)
}
