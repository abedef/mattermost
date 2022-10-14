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

func Listen(address string) {
	for {
		webSocketClient, err := model.NewWebSocketClient4(address, getClient().AuthToken)
		if err != nil {
			log.Print(err)
		}
		log.Print("Connected to WS")
		webSocketClient.Listen()

		for resp := range webSocketClient.EventChannel {
			log.Print(resp)
		}
	}
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

func createChannel(name string, tname string) {
	team, resp := getClient().GetTeamByName(tname, "")
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}
	channel := &model.Channel{}
	channel.Name = name
	channel.DisplayName = name
	channel.Purpose = ""
	channel.Type = model.CHANNEL_PRIVATE
	channel.TeamId = team.Id
	rchannel, resp := getClient().CreateChannel(channel)
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}

	getClient().AddChannelMember(rchannel.Id, "1dm71jzndf8jtedqbbdx7uc3co")
	getClient().AddChannelMember(rchannel.Id, "n4hw1yqaq7rzpjgitxq7uqj9ke")
}

func WritePostToChannel(channel string, team string, text string) {
	rchannel, resp := getClient().GetChannelByNameForTeamName(channel, team, "")
	if resp.Error != nil {
		createChannel(channel, team)
		rchannel, resp = getClient().GetChannelByNameForTeamName(channel, team, "")
		if resp.Error != nil {
			log.Fatal(resp.Error)
		}
	}

	post := model.Post{
		ChannelId: rchannel.Id,
		Message:   text,
	}

	_, response := getClient().CreatePost(&post)

	log.Print(response)
}
