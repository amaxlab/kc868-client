package main

import (
	"fmt"
	"github.com/amaxlab/go-lib/log"
)

func main() {
	config := NewConfiguration()
	client := NewKC868Client(config.KC868Config.Host, config.KC868Config.Port)
	client.connect()
	defer client.disconnect()

	if config.Debug {
		log.Debug.Enable()
	}

	log.Info.Printf(fmt.Sprintf("Listen: %d port", config.Port))

	webServer := NewWebServer(config.Port, client)

	err := webServer.start()
	if err != nil {
		log.Error.Fatal(err)
	}
}
