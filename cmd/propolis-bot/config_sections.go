package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/logthis"
)

const (
	centralRole    = "central"
	nodeRole       = "node"
	defaultChannel = "#propolis-announce"
	centralBot     = "bbb"
)

type ConfigGeneral struct {
	LogLevel      int    `yaml:"log_level"`
	PrivateBinURL string `yaml:"privatebin_url"`
	PtpImgKey     string `yaml:"ptpimg_key"`
}

func (cg *ConfigGeneral) check() error {
	if cg.PrivateBinURL == "" {
		return errors.New("privatebin URL required")
	}
	if cg.PtpImgKey == "" {
		return errors.New("PTPIMG API key required")
	}
	return logthis.CheckLevel(cg.LogLevel)
}

// String representation for ConfigGeneral.
func (cg *ConfigGeneral) String() string {
	txt := "General configuration:\n"
	txt += "\tLog level: " + strconv.Itoa(cg.LogLevel) + "\n"
	txt += "\tPrivateBin URL: " + cg.PrivateBinURL + "\n"
	txt += "\tPTPImg API key: " + cg.PtpImgKey + "\n"
	return txt
}

type ConfigIrc struct {
	Server           string `yaml:"server"`
	Username         string `yaml:"username"`
	Key              string `yaml:"key"`
	UseSSL           bool   `yaml:"ssl"`
	NickServPassword string `yaml:"nickserv_password"`
	BotName          string `yaml:"bot_name"`
	GateKeeper       string `yaml:"gatekeeper"`
	Channel          string `yaml:"channel"`
	Role             string `yaml:"role"`
	CentralBot       string `yaml:"central_bot"`
}

func (ci *ConfigIrc) check() error {
	// TODO keep the whole section optional
	if ci.Role == "" {
		ci.Role = nodeRole
	} else if ci.Role != nodeRole && ci.Role != centralRole {
		return errors.New("invalid role")
	}
	if ci.Channel == "" {
		ci.Channel = defaultChannel
	}
	if ci.CentralBot == "" {
		ci.CentralBot = centralBot
	}
	return nil
}

func (ci *ConfigIrc) String() string {
	txt := "IRC configuration:\n"
	txt += "\tServer: " + ci.Server + "\n"
	txt += "\tIRC key: " + ci.Key + "\n"
	txt += "\tUse SSL: " + fmt.Sprintf("%v", ci.UseSSL) + "\n"
	txt += "\tNickServ password: " + ci.NickServPassword + "\n"
	txt += "\tUser Name: " + ci.Username + "\n"
	txt += "\tBot Name: " + ci.BotName + "\n"
	txt += "\tGateKeeper: " + ci.GateKeeper + "\n"
	txt += "\tChannel: " + ci.Channel + "\n"
	txt += "\tRole: " + ci.Role + "\n"
	txt += "\tCentral Bot: " + ci.CentralBot + "\n"
	return txt
}

type ConfigTracker struct {
	Site                 string         `yaml:"site"`
	Token                string         `yaml:"token"`
	Port                 int            `yaml:"port"`
	SessionCookie        string         `yaml:"session_cookie"`
	APIKey               string         `yaml:"api_key"`
	URL                  string         `yaml:"tracker_url"`
	BlacklistedUploaders []string       `yaml:"blacklisted_uploaders"`
	ExcludedTags         []string       `yaml:"excluded_tags"`
	Users                []string       `yaml:"whitelisted_users"`
	WhitelistedUsers     map[string]int `yaml:"-"`
}

func (cv *ConfigTracker) check() error {
	if cv.Site != "" && cv.Token == "" {
		return errors.New("missing token")
	}
	if cv.Site == "" && cv.Token != "" {
		return errors.New("missing site name")
	}
	if cv.Site != "" && cv.Port == 0 {
		return errors.New("missing port number")
	}
	if cv.SessionCookie == "" || cv.APIKey == "" {
		return errors.New("both an API key and a session cookie are required")
	}

	// TODO more checks

	// parsing the whitelisted users configuration
	cv.WhitelistedUsers = make(map[string]int)
	for _, u := range cv.Users {
		parts := strings.Split(u, ":")
		if len(parts) == 2 {
			number, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
			cv.WhitelistedUsers[parts[0]] = number
		}
	}

	return nil
}

func (cv *ConfigTracker) String() string {
	txt := "Tracker configuration:\n"
	txt += "\tSite: " + cv.Site + "\n"
	txt += "\tToken: " + cv.Token + "\n"
	txt += "\tPort: " + fmt.Sprintf("%d", cv.Port) + "\n"
	txt += "\tTracker URL: " + cv.URL + "\n"
	txt += "\tTracker API Key: " + cv.APIKey + "\n"
	txt += "\tTracker Session Cookie: " + cv.SessionCookie + "\n"
	txt += "\tBlacklisted uploaders: " + strings.Join(cv.BlacklistedUploaders, ", ") + "\n"
	txt += "\tBlacklisted tags: " + strings.Join(cv.ExcludedTags, ", ") + "\n"
	var whitelisted []string
	for k, v := range cv.WhitelistedUsers {
		if v == -1 {
			whitelisted = append(whitelisted, fmt.Sprintf("%s: infinite!", k))
		} else {
			whitelisted = append(whitelisted, fmt.Sprintf("%s: %d / day", k, v))
		}
	}
	txt += "\tWhitelisted Users: " + strings.Join(whitelisted, ", ") + "\n"
	return txt
}
