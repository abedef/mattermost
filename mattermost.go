package mattermost

import (
	"log"

	"github.com/mattermost/mattermost-server/model"
)

var client *model.Client4

var Url, Token string

var hasBeenConfigured = false

func Configure(url string, token string) {
	Url = url
	Token = token

	hasBeenConfigured = true
}

func getClient() *model.Client4 {
	if client != nil {
		return client
	}

	if !hasBeenConfigured {
		log.Fatalln("failed to initialize client: did not call Configure() to supply configuration values")
	}

	client := model.NewAPIv4Client(Url)
	client.SetOAuthToken(Token)

	return client
}

func writePostToChannel(channel string, text string) {
	post := model.Post{
		ChannelId: channel,
		Message:   text,
	}

	_, response := getClient().CreatePost(&post)

	log.Print(response)
}

func WriteProdPost(text string) {
	writePostToChannel("qz1wyhymmfd5imeaerh83ydpnh", text)
}

func WritePostDev(text string) {
	writePostToChannel("8wa1jsbiui859ni8qicgo9314w", text)
}
