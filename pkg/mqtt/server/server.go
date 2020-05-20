package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/naspinall/Hive/pkg/mqtt/packets"
)

func NewMQTTBroker() MQTT {
	return MQTT{
		SubscriptionHandlers: make(map[string]SubscribeHandler),
		PublishHandlers:      make(map[string]PublishHandler),
		// Default Auth handler
		AuthHandler: func(b []byte) (bool, error) {
			return true, nil
		},
	}
}

type PublishHandler func(*packets.PublishPacket)
type SubscribeHandler func(*packets.SubscribePacket, *Connection)

type MQTT struct {
	SubscriptionHandlers map[string]SubscribeHandler
	PublishHandlers      map[string]PublishHandler
	AuthHandler          func(b []byte) (bool, error)
}

func (mqtt *MQTT) Listen(host string, port string) {

	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go mqtt.HandleNewConn(conn)
	}
}

func (mqtt *MQTT) Publish(topic string, handler PublishHandler) {
	mqtt.PublishHandlers[topic] = handler
}

func (mqtt *MQTT) Subscribe(topic string, handler SubscribeHandler) {
	mqtt.SubscriptionHandlers[topic] = handler
}

func (mqtt *MQTT) HandleNewConn(conn net.Conn) {
	b := make([]byte, 4096)
	n, err := conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d byte from new connection", n)

	// Decoding packets
	fh, err := packets.NewFixedHeader(b)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Checking if first packet sent is a connect packet
	if fh.Type != packets.CONNECT {
		log.Println("Bad packet send")
		conn.Close()
		return
	}

	// Creating session.
	c, err := mqtt.NewConnection(conn, 1)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Sending accepted response
	cb, err := packets.Accepted().Encode()
	if err != nil {
		log.Println("Cannot encode Accepted packet")
		conn.Close()
		return
	}
	n, err = c.Conn.Write(cb)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Handling the connection
	go mqtt.HandleConnection(c)
}

func (mqtt *MQTT) NewConnection(conn net.Conn, deviceID uint) (*Connection, error) {
	c := &Connection{
		Session: &Session{
			SessionID: uuid.New().String(),
			DeviceID:  deviceID,
		},
		Conn: conn,
	}

	return c, nil
}

func (mqtt *MQTT) HandleConnection(c *Connection) {
	for {
		fh, b, err := c.ProcessFixedHeader()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch fh.Type {
		case packets.PUBLISH:
			pp, err := packets.NewPublishPacket(fh, b)
			if err != nil {
				log.Println(err)
			}
			switch pp.FixedHeader.Flags.QoS {
			case 1:
				b := make([]byte, 4)
				if _, err := packets.Acknowledge(pp.PacketIdentifier).Encode(b); err != nil {
					log.Println("Back Ack packet encoding")
				}
				c.Conn.Write(b)
			case 2:
				b := make([]byte, 4)
				if _, err := packets.Received(pp.PacketIdentifier).Encode(b); err != nil {
					log.Println("Back Ack packet encoding")
				}
				c.Conn.Write(b)
				rc := make(chan uint16)
				timeOut := time.NewTimer(500 * time.Microsecond)
				go c.PublishQos(rc)
				select {
				case pi := <-rc:
					if _, err := packets.Complete(pi).Encode(b); err != nil {
						log.Println("Back Ack packet encoding")
					}
					c.Conn.Write(b)
				case <-timeOut.C:
					fmt.Println("Timed out")
					continue
				}
			}
			// Support for wildcards and multilevel coming soon
			if handler, ok := mqtt.PublishHandlers[pp.TopicName]; ok == true {
				handler(pp)
			}

		case packets.SUBSCRIBE:
			sp, err := packets.NewSubscribePacket(fh, b)
			if err != nil {
				log.Println(err)
			}
			// Support for wildcards and multilevel coming soon
			for _, topic := range sp.Topics {
				if handler, ok := mqtt.SubscriptionHandlers[topic.Topic]; ok == true {
					go handler(sp, c)
				}
			}
		case packets.UNSUBSCRIBE:
			continue
		case packets.PINGREQ:
			log.Println("<-- PING")
			pr, err := packets.PingResponse().Encode()
			if err != nil {
				log.Println(err)
			}
			c.Conn.Write(pr)
			log.Println("PONG -->")
		case packets.DISCONNECT:
			// Removing the session from the database, disconnecting.
			c.Conn.Close()
			break
		default:
			continue
		}
	}
}

func (c *Connection) PublishQos(rc chan uint16) {
	b := make([]byte, 4)
	for {
		_, err := c.Conn.Read(b)
		if err != nil {
			log.Println("Connection read error")
			return
		}
		fh, err := packets.NewFixedHeader(b)
		if fh.Type == 6 {
			pr, err := packets.NewPublishQoSPacket(fh, b)
			if err != nil {
				log.Println("Bad publish QoS Packet Provided")
			}
			rc <- pr.PacketIdentifier.PacketIdentifier
		}
	}
}

func (c *Connection) ProcessFixedHeader() (fh *packets.FixedHeader, b []byte, err error) {
	b = make([]byte, 2)
	_, err = c.Conn.Read(b)
	fh, err = packets.NewFixedHeader(b)
	if err != nil {
		log.Println(err)
	}

	// Reading only the length of the packet
	b = make([]byte, fh.RemaningLength)
	_, err = c.Conn.Read(b)
	return
}

// TODO
// Improve Networking
// Use Context API for connections
// Work on retain for multiplexing to other subscribers
