package client

import (
	"encoding/json"
	"fmt"

	"github.com/RadiumByte/LabYoutubeChatbot/app"
	"github.com/valyala/fasthttp"
)

// SelectCameraJSON represents transport data for camera switching
type SelectCameraJSON struct {
	CameraName string `json:"name"`
}

// ServerClient represents data for connection to Stream Server
type ServerClient struct {
	Client     *fasthttp.Client
	Request    *fasthttp.Request
	Response   *fasthttp.Response
	ServerIP   string
	ServerPort string
}

// GetCameras receives list of all available cameras from Stream Server
func (c *ServerClient) GetCameras() []app.CameraData {
	c.Request.Header.SetMethod("GET")

	url := "http://" + c.ServerIP + c.ServerPort + "/get-cameras"

	c.Request.SetRequestURI(url)
	err := c.Client.Do(c.Request, c.Response)

	if err != nil {
		fmt.Println("Client: GetCameras failed to make a request.")
		return nil
	}

	payload := c.Response.Body()
	var dataJSON map[string]interface{}

	if err := json.Unmarshal(payload, &dataJSON); err != nil {
		fmt.Println("Client: Server returned bad data for GetCameras")
		return nil
	}

	var types []interface{}
	var names []interface{}
	var cameras []app.CameraData

	types = dataJSON["types"].([]interface{})
	names = dataJSON["names"].([]interface{})

	for i := 0; i < len(types); i++ {
		current := app.CameraData{
			Name: names[i].(string),
			Type: types[i].(int)}
		cameras = append(cameras, current)
	}

	return cameras
}

// GetActive gets one active (broadcasting) camera at this moment
func (c *ServerClient) GetActive() error {
	c.Request.Header.SetMethod("GET")

	url := "http://" + c.ServerIP + c.ServerPort + "/get-active"

	c.Request.SetRequestURI(url)
	c.Client.Do(c.Request, c.Response)

	return nil
}

// SelectCamera makes specified camera active, switching the broadcast
func (c *ServerClient) SelectCamera(name string) {
	c.Request.Header.SetMethod("POST")
	c.Request.Header.SetContentType("application/json")

	url := "http://" + c.ServerIP + c.ServerPort + "/select-camera"
	c.Request.SetRequestURI(url)

	toEncode := &SelectCameraJSON{
		CameraName: name}

	payload, _ := json.Marshal(toEncode)

	c.Request.SetBody(payload)

	c.Client.Do(c.Request, c.Response)
}

// NewServerClient constructs object of ServerClient
func NewServerClient(ip string, port string) (*ServerClient, error) {
	res := &ServerClient{}
	res.Client = &fasthttp.Client{}
	res.Request = fasthttp.AcquireRequest()
	res.Response = fasthttp.AcquireResponse()
	res.ServerPort = port
	res.ServerIP = ip

	return res, nil
}
