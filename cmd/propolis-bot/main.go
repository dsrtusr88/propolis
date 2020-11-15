package main

import (
	"fmt"
)

const (
	fullName = "propolisBot"
)

var Version = "dev"

func userAgent() string {
	return fullName + "/" + Version
}

func main() {
	var err error
	config, err = NewConfig(DefaultConfigurationFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	if !config.IrcConfigured {
		fmt.Println("irc section is not configured")
		return
	}

	// launching the websocket server
	go Server()

	// only join a channel if central bot
	var channels []string
	if config.IRC.Role == centralRole {
		channels = []string{config.IRC.Channel}
	}

	// launching the bot
	Run(&BotConfig{
		Server:     config.IRC.Server,
		Channels:   channels,
		Nick:       config.IRC.Username,
		User:       config.IRC.BotName,
		RealName:   "BUSY BEE BOT",
		Password:   config.IRC.NickServPassword,
		Key:        config.IRC.Key,
		GateKeeper: config.IRC.GateKeeper,
		UseTLS:     true,
		//	Debug:      true,
	})
}
