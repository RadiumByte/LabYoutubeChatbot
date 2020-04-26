package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// YouTubeChatbot is an interface for application's core
type YouTubeChatbot interface {
}

// StreamServerClient is an interface for calling Stream Server from chatbot
type StreamServerClient interface {
	GetCameras() error
	GetActive() error
	SelectCamera(name string)
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	server StreamServerClient
}

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	fmt.Println(tokenCacheDir)
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func getBroadcastID(service *youtube.Service, part string) string {
	call := service.LiveBroadcasts.List(part)
	call = call.Mine(true)

	response, err := call.Do()
	handleError(err, "")

	fmt.Println(fmt.Sprintf("This broadcast's ID is %s. It's live chat ID is '%s', "+
		"and it's description: %s.",
		response.Items[0].Id,
		response.Items[0].Snippet.LiveChatId,
		response.Items[0].Snippet.Description))

	return response.Items[0].Snippet.LiveChatId
}

func getMessages(service *youtube.Service, part string, chatID string, pageToken string) (string, []*youtube.LiveChatMessage) {
	call := service.LiveChatMessages.List(chatID, part)

	if pageToken != "" {
		call.PageToken(pageToken)
	}

	response, err := call.Do()
	handleError(err, "")

	for i := 0; i < len(response.Items); i++ {
		fmt.Printf("Message: %s\n", response.Items[i].Snippet.TextMessageDetails.MessageText)
	}
	fmt.Printf("Page token: %s\n", response.NextPageToken)

	return response.NextPageToken, response.Items
}

func sendMessage(service *youtube.Service, part string, chatID string, message string) {
	toSend := &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: message,
			},
			LiveChatId: chatID,
			Type:       "textMessageEvent",
		},
	}

	call := service.LiveChatMessages.Insert(part, toSend)

	response, err := call.Do()
	fmt.Println(err)
	fmt.Println(response)
	handleError(err, "")
}

func getViewers(service *youtube.Service, part string) uint64 {
	call := service.LiveBroadcasts.List(part)
	call = call.Mine(true)

	response, err := call.Do()
	handleError(err, "")
	return response.Items[0].Statistics.ConcurrentViewers
}

// Start runs all connecting and parsing process
func (a *Application) Start() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	chatID := getBroadcastID(service, "snippet")

	pageToken := ""

	var items []*youtube.LiveChatMessage

	var lastViewers uint64
	var currentViewers uint64

	lastViewers = 0
	currentViewers = 0

	for true {

		lastViewers = currentViewers
		currentViewers = getViewers(service, "statistics")

		if currentViewers > lastViewers {
			sendMessage(service, "snippet", chatID, "Добро пожаловать в чат Лаборатории ИИ и робототехники мехмата ЮФУ!")
			sendMessage(service, "snippet", chatID, "Для управления камерами вы можете воспользоваться следующими командами:")
			sendMessage(service, "snippet", chatID, "'Список камер'")
			sendMessage(service, "snippet", chatID, "'Выбрать камеру <название>'")
			sendMessage(service, "snippet", chatID, "'Активная камера'")
		}

		pageToken, items = getMessages(service, "snippet", chatID, pageToken)

		for i := len(items) - 1; i >= 0; i-- {
			currentMessage := items[i].Snippet.TextMessageDetails.MessageText
			if currentMessage == "Список камер" {
				a.server.GetCameras()
				break
			} else if currentMessage == "Активная камера" {
				a.server.GetActive()
			} else if strings.Index(currentMessage, "Выбрать камеру") != -1 {

			}
		}

		time.Sleep(5 * time.Second)
	}
}

// NewApplication constructs Application
func NewApplication(serverClient StreamServerClient) *Application {
	res := &Application{}
	res.server = serverClient

	res.Start()
	return res
}
