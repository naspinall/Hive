package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/naspinall/Hive/pkg/config"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/naspinall/Hive/pkg/mqtt/client"
	"github.com/naspinall/Hive/pkg/mqtt/models"
	"github.com/naspinall/Hive/pkg/mqtt/packets"
)

func NewMQTTBroker() MQTT {
	pc := config.LoadFromEnvironment()

	services, err := models.NewServices(
		models.WithGorm("postgres", pc.ConnectionInfo()),
		models.WithRetain(),
		models.WithSession(),
		models.WithWill(),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = services.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	return MQTT{
		// Default Auth handler
		AuthHandler: func(b []byte) (bool, error) {
			return true, nil
		},
		Subscriptions: make(map[string][]*Connection),
		Services:      *services,
	}
}

type PublishHandler func(*packets.PublishPacket)
type SubscribeHandler func(*packets.SubscribePacket, *Connection)

type MQTT struct {
	models.Services
	Subscriptions map[string][]*Connection
	AuthHandler   func(b []byte) (bool, error)
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

func (mqtt *MQTT) HandlePublish(pp *packets.PublishPacket) error {

	sessions, ok := mqtt.Subscriptions[pp.TopicName]

	if ok {
		for _, session := range sessions {
			fmt.Printf("Sending for topic %s", pp.TopicName)
			err := client.Publish(pp, session.Conn)
			// One error shouldn't break all of the publishes.
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (mqtt *MQTT) HandleSubscribe(pp *packets.SubscribePacket, c *Connection) {

	for _, topic := range pp.Topics {
		// If subscription alreay exists we'll add to the curernt list of connections
		current, ok := mqtt.Subscriptions[topic.Topic]
		if ok {
			current = append(current, c)
			continue
		}
		fmt.Printf("%+v", mqtt.Subscriptions[topic.Topic])
		mqtt.Subscriptions[topic.Topic] = []*Connection{c}
	}

	client.SubAck(pp.PacketIdentifier.PacketIdentifier, 0, c.Conn)
}

func (mqtt *MQTT) HandleNewConn(conn net.Conn) {
	p, err := packets.FromReader(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// Checking if first packet sent is a connect packet
	if p.Type != packets.CONNECT {
		log.Println("Inital packet is not a connect packet")
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
	log.Println("Sending Accept")
	if err != nil {
		log.Println("Cannot encode Accepted packet")
		conn.Close()
		return
	}
	_, err = c.Conn.Write(cb)
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
		ClientID: uuid.New().String(),
		Conn:     conn,
	}

	return c, nil
}

func (mqtt *MQTT) HandleConnection(c *Connection) {
	for {
		p, err := packets.FromReader(c.Conn)
		if err != nil {
			fmt.Println("Bad packet read")
			c.Conn.Close()
			return
		}

		switch p.Type {
		case packets.PUBLISH:
			pp, err := packets.NewPublishPacket(p)
			if err != nil {
				log.Println(err)
				break
			}
			mqtt.HandlePublish(pp)
			switch pp.Flags.QoS {
			case 1:
				b, err := packets.Acknowledge(pp.PacketIdentifier).Encode()
				if err != nil {
					log.Println("Back Ack packet encoding")
					break
				}
				c.Conn.Write(b)
			case 2:
				b, err := packets.Received(pp.PacketIdentifier).Encode()
				if err != nil {
					log.Println("Back Ack packet encoding")
					break
				}
				c.Conn.Write(b)
				rc := make(chan uint16)
				timeOut := time.NewTimer(500 * time.Microsecond)

				select {
				case pi := <-rc:
					b, err := packets.Complete(pi).Encode()
					if err != nil {
						log.Println("Back Ack packet encoding")
					}
					c.Conn.Write(b)
				case <-timeOut.C:
					fmt.Println("Timed out")
					continue
				}
			}

		case packets.SUBSCRIBE:
			sp, err := packets.NewSubscribePacket(p)
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println("Subsribe me!")
			mqtt.HandleSubscribe(sp, c)
		case packets.UNSUBSCRIBE:
			continue
		case packets.PINGREQ:
			log.Println("<-- PING")
			pr, err := packets.PingResponse().Encode()
			if err != nil {
				log.Println(err)
				break
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

// func (c *Connection) PublishQos(rc chan uint16) {
// 	b := make([]byte, 4)
// 	for {
// 		_, err := c.Conn.Read(b)
// 		if err != nil {
// 			log.Println("Connection read error")
// 			return
// 		}
// 		fh, err := packets.NewFixedHeader(b)
// 		if fh.Type == 6 {
// 			pr, err := packets.NewPublishQoSPacket(fh, b)
// 			if err != nil {
// 				log.Println("Bad publish QoS Packet Provided")
// 			}
// 			rc <- pr.PacketIdentifier.PacketIdentifier
// 		}
// 	}
// }

// TODO
// Improve Networking
// Use Context API for connections
// Work on retain for multiplexing to other subscribers
