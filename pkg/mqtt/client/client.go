package client

import (
	"fmt"
	"net"

	"github.com/google/uuid"

	"github.com/naspinall/Hive/pkg/mqtt/packets"
)

type Client struct {
	conn net.Conn
}

type ClientOptions struct {
	ConnectPacket packets.ConnectPacket
	ClientID      string
	Address       string
}

func (co *ClientOptions) NewClient() (*Client, error) {

	conn, err := net.Dial("tcp", co.Address)
	b, err := co.ConnectPacket.Encode()
	_, err = conn.Write(b)
	if err != nil {
		return nil, err
	}
	rb := make([]byte, 4)
	_, err = conn.Read(rb)
	if err != nil {
		return nil, err
	}

	fmt.Println(rb)

	return &Client{conn}, nil
}

func (co *ClientOptions) SetUsername(p string) *ClientOptions {
	co.ConnectPacket.UsernameFlag = true
	co.ConnectPacket.Username = p
	return co
}

func (co *ClientOptions) SetPassword(p []byte) *ClientOptions {
	co.ConnectPacket.PasswordFlag = true
	co.ConnectPacket.Password = p
	return co
}

func (co *ClientOptions) SetWillPayload(p []byte) *ClientOptions {
	co.ConnectPacket.WillPayload = p
	return co
}

func (co *ClientOptions) SetWillRetain(p bool) *ClientOptions {
	co.ConnectPacket.WillRetainFlag = p
	return co
}

func (co *ClientOptions) SetWillQoS(p uint8) *ClientOptions {
	co.ConnectPacket.WillQoSFlag = p
	return co
}

func (co *ClientOptions) SetCleanSession(p bool) *ClientOptions {
	co.ConnectPacket.CleanSessionFlag = p
	return co
}

func (co *ClientOptions) SetKeepAlive(p uint16) *ClientOptions {
	co.ConnectPacket.KeepAlive = p
	return co
}

func (co *ClientOptions) SetAddress(p string) *ClientOptions {
	co.Address = p
	return co
}

func (co *ClientOptions) SetClientID(p string) *ClientOptions {
	co.ClientID = p
	return co
}

func NewClientOptions() *ClientOptions {
	clientID := uuid.New().String()
	return &ClientOptions{
		ConnectPacket: packets.ConnectPacket{
			FixedHeader: &packets.FixedHeader{
				Type: packets.CONNECT,
			},
			ProtocolName:    "MQTT",
			ProtocolVersion: 0x04,
			KeepAlive:       uint16(60),

			//Payload properties
			ClientID: clientID,
		},
	}
}

func (c *Client) Publish(topic string, payload []byte) {

}
