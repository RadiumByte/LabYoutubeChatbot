package main

import (
	"github.com/RadiumByte/LabYoutubeChatbot/app"
	"github.com/RadiumByte/LabYoutubeChatbot/client"

	"fmt"
)

func main() {
	fmt.Println("YouTube Live Chatbot is preparing to start...")

	serverClient, err := client.NewServerClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	app.NewApplication(serverClient)
}
