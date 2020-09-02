package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

// User ...
type User struct {
	discord *discordgo.User
	spotify *spotify.Client
}

// NewUser ...
func NewUser(discord *discordgo.User, spotify *spotify.Client) *User {
	user := new(User)
	user.discord = discord
	user.spotify = spotify

	return user
}
