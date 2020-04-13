package client

// ServerClient represents data for connection to Stream Server
type ServerClient struct {
}

// NewServerClient constructs object of ServerClient
func NewServerClient() (*ServerClient, error) {
	res := &ServerClient{}

	return res, nil
}
