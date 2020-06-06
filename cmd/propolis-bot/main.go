package main

import (
	"fmt"
)

const (
	errorMsg = "error"
)

func main() {
	var err error
	config, err = NewConfig(DefaultConfigurationFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !config.IrcConfigured {
		fmt.Println("irc section is not configured")
		return
	}

	// launching the websocket server
	go Server()

	// launching the bot
	Run(&BotConfig{
		Server:     config.IRC.Server,
		Channels:   []string{config.IRC.Channel},
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
