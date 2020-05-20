package main

import "github.com/naspinall/Hive/pkg/mqtt/client"

func main() {
	opts := client.NewClientOptions().SetAddress("localhost:8080")
	opts.NewClient()
}
