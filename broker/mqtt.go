package main

import (
	"github.com/naspinall/Hive/pkg/mqtt/packets"
	"github.com/naspinall/Hive/pkg/mqtt/server"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "hive"
	dbname   = "hive"
)

func main() {
	mqtt := server.NewMQTTBroker(nil)
	mqtt.Publish("hello", func(pp *packets.PublishPacket) {
		//fmt.Println(string(pp.Payload))
	})
	mqtt.Listen("localhost", "8080")
}
