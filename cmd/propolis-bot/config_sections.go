package main

import (
	"fmt"
	"strconv"

	"gitlab.com/catastrophic/assistance/logthis"
)

type ConfigGeneral struct {
	LogLevel int `yaml:"log_level"`
}

func (cg *ConfigGeneral) check() error {
	return logthis.CheckLevel(cg.LogLevel)
}

// String representation for ConfigGeneral.
func (cg *ConfigGeneral) String() string {
	txt := "General configuration:\n"
	txt += "\tLog level: " + strconv.Itoa(cg.LogLevel) + "\n"
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
}

func (ci *ConfigIrc) check() error {
	// TODO keep the whole section optional
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
	return txt
}
