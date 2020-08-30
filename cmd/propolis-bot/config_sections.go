package main

import (
	"fmt"
	"strconv"

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
