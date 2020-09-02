package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

var (
	// Parties ...
	Parties *PartyManager
	wg      = new(sync.WaitGroup)
	// Auth ...
	Auth spotify.Authenticator
)

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
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

	ctx := NewContext(session, guild, channel, message.Author, message, Parties, Auth)

	if message.Content == "!join" {
		JoinCommand(ctx)
	}

	if message.Content == "!list" {
		ListCommand(ctx)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wg.Add(2)
	Parties = NewPartyManager()
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
