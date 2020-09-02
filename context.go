package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

// Context ...
type Context struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Message *discordgo.MessageCreate
	Parties *PartyManager
	Auth    spotify.Authenticator
}

// NewContext ..
func NewContext(session *discordgo.Session, guild *discordgo.Guild, channel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, parties *PartyManager, auth spotify.Authenticator) *Context {

	ctx := new(Context)
	ctx.Session = session
	ctx.Guild = guild
	ctx.Channel = channel
	ctx.User = user
	ctx.Message = message
	ctx.Parties = parties
	ctx.Auth = auth

	return ctx
}

// Reply ..
func (ctx Context) Reply(content string) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, content)

	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
	return err
}
