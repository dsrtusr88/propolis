package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-chat-bot/bot"
)

func hello(command *bot.Cmd) (string, error) {
	return fmt.Sprintf("Why, hello there %s! You can indeed call me %s!", command.User.Nick, command.RawArgs), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	bot.RegisterCommand(
		"hello",
		"Does nothing valuable.",
		"bot",
		hello)
}
