package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	outputFormat = "\u0002\u000371%s>\u000F %s"
)

func ParseNodeMessage(sender, message string) (string, string, string, error) {
	parts := strings.Split(message, "|")
	if len(parts) == 2 {
		if strings.Contains(message, "Error: ") {
			return parts[0], parts[1], fmt.Sprintf(outputFormat, sender, message), nil
		}
	}
	if len(parts) != 5 {
		return "", "", "", errors.New("unrecognized input from " + sender + ": " + message)
	}
	// TODO: other checks, such as ID == int
	if !strings.Contains(parts[3], config.General.PrivateBinURL) {
		return "", "", "", errors.New("unrecognized input from " + sender + ": " + message)
	}
	return parts[0], parts[1], fmt.Sprintf(outputFormat, sender, message), nil
}
