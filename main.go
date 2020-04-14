package main

import (
	"github.com/RadiumByte/LabYoutubeChatbot/app"
	"github.com/RadiumByte/LabYoutubeChatbot/client"

	"fmt"
)

func main() {
	fmt.Println("YouTube Live Chatbot is preparing to start...")
	serverIP := "localhost"
	serverPort := ":8081"

	serverClient, err := client.NewServerClient(serverIP, serverPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.NewApplication(serverClient)
}
