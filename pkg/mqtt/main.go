package main

import (
	"fmt"
	"log"

	"github.com/naspinall/Hive/mqtt/packets"
	"github.com/naspinall/Hive/mqtt/server"

	"github.com/jinzhu/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "hive"
	dbname   = "hive"
)

func main() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&server.Session{}).Error; err != nil {
		log.Fatal(err)
	}

	mqtt := server.NewMQTTBroker(db)
	mqtt.Publish("hello", func(pp *packets.PublishPacket, c *server.Connection) {
		fmt.Println("Recieved publish!")
		c.Conn.Close()
	})
	mqtt.Listen("localhost", "8080")
}
