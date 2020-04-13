package main

import (
	"github.com/RadiumByte/LabYoutubeChatbot/app"
	"github.com/RadiumByte/LabYoutubeChatbot/client"

	"fmt"
)

func run(errc chan<- error) {
	serverClient, err := client.NewServerClient()
	if err != nil {
		errc <- err
		return
	}

	application := app.NewApplication(serverClient, errc)
}

func main() {
	fmt.Println("Youtube Live Chatbot is preparing to start...")

	errc := make(chan error)
	go run(errc)
	if err := <-errc; err != nil {
		fmt.Print("Error occured: ")
		fmt.Println(err)
	}
}
