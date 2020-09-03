package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

var (
	Parties    *PartyManager
	wg         = new(sync.WaitGroup)
	Auth       spotify.Authenticator
	CmdHandler *CommandHandler
	PREFIX     string
)

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == session.State.User.ID {
		return
	}
	content := message.Content
	if len(content) <= len(PREFIX) || content[:len(PREFIX)] != PREFIX {
		return
	}

	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}

	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, err := CmdHandler.Get(name)
	if err != nil {
		return
	}

	channel, err := session.State.Channel(message.ChannelID)

	if err != nil {
		fmt.Println("Error getting channel: ", err)
		return
	}

	guild, err := session.State.Guild(message.GuildID)

	if err != nil {
		fmt.Println("Error getting guild: ", err)
		return
	}

	ctx := NewContext(session, guild, channel, user, message, Parties, Auth, args[1:])
	c := *command
	c(ctx)
}

func main() {
	err := godotenv.Load()
	PREFIX = os.Getenv("PREFIX")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wg.Add(2)
	Parties = NewPartyManager()
	CmdHandler = NewCommandHandler()
	Auth = GetSpotifyAuth()
	go InitDiscord()
	go InitAuthServer()
	wg.Wait()
}

// GetSpotifyAuth ...
func GetSpotifyAuth() spotify.Authenticator {
	clientID := os.Getenv("SPOTIFY_ID")
	secretKey := os.Getenv("SPOTIFY_SECRET")
	url := os.Getenv("URL")

	auth := spotify.NewAuthenticator(url, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, secretKey)

	return auth
}
