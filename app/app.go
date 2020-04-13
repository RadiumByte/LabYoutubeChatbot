package app

// YouTubeChatbot is an interface for application's core
type YouTubeChatbot interface {
}

// StreamServerClient is an interface for calling Stream Server from chatbot
type StreamServerClient interface {
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	server StreamServerClient
}

// Start runs all connecting and parsing process
func (a *Application) Start() {

}

// NewApplication constructs Application
func NewApplication(serverClient StreamServerClient) *Application {
	res := &Application{}
	res.server = serverClient

	res.Start()
	return res
}
