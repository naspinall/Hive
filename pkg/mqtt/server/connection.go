package server

import (
	"net"
)

type Connection struct {
	Session *Session
	Conn    net.Conn
}
