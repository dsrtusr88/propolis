package propolis

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"gitlab.com/catastrophic/assistance/logthis"
)

const (
	WebsocketPort    = 8335
	HandshakeCommand = "hello"
	LogCommand       = "log"
)

const (
	responseInfo = iota
	responseError
)

// IncomingJSON from here to the server.
type IncomingJSON struct {
	Command string
	Args    []string
}

// OutgoingJSON from the server to here.
type OutgoingJSON struct {
	Status  int
	Message string
}

func SendToWebsocket(lines []string) error {
	//u := url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1:%d", WebsocketPort), Path: "/ws"}
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("172.17.0.1:%d", WebsocketPort), Path: "/ws"}
	logthis.Info("connecting to "+u.String(), logthis.VERBOSESTEST)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	var i IncomingJSON
	var j OutgoingJSON

	// performing handshake
	i.Command = HandshakeCommand
	if err := c.WriteJSON(i); err != nil {
		return err
	}
	if err := c.ReadJSON(&j); err != nil {
		return err
	}
	if j.Status != responseInfo || j.Message != HandshakeCommand {
		return errors.New("unsuccessful handshake")
	} else {
		logthis.Info("handshake completed with "+u.String(), logthis.VERBOSESTEST)
	}

	// sending the log
	i.Command = LogCommand
	i.Args = lines
	if err := c.WriteJSON(i); err != nil {
		return err
	}
	if err := c.ReadJSON(&j); err != nil {
		return err
	}
	if j.Status != responseInfo || j.Message != "OK" {
		return errors.New("could not send lines to propolis-bot")
	} else {
		logthis.Info("lines sent to propolis-bot at "+u.String(), logthis.VERBOSESTEST)
	}
	return nil
}
