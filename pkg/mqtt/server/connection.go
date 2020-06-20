package server

import (
	"net"
)

type Connection struct {
	ClientID string
	Conn     net.Conn
}
