package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	fmt.Println("+ Testing Config...")
	check := assert.New(t)

	c := &Config{}
	err := c.Load("../../test/config.yaml")
	check.Nil(err)

	// general
	fmt.Println("Checking general")
	check.Equal(3, c.General.LogLevel)

	// irc
	check.True(c.IrcConfigured)
	check.Equal("irc.server.net:6697", c.IRC.Server)
	check.Equal("kkeeyy", c.IRC.Key)
	check.True(c.IRC.UseSSL)
	check.Equal("something", c.IRC.NickServPassword)
	check.Equal("user", c.IRC.Username)
	check.Equal("mybot", c.IRC.BotName)
	check.Equal("Bee", c.IRC.GateKeeper)
	check.Equal(defaultChannel, c.IRC.Channel)
	check.Equal(nodeRole, c.IRC.Role)
	check.Equal(centralBot, c.IRC.CentralBot)

	fmt.Println(c.String())
}
