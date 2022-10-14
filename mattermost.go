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
	_, resp = getClient().CreateChannel(channel)
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}
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
