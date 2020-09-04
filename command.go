package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

type (
	Command         func(*Context)
	CommandWithHelp struct {
		command Command
		help    string
	}
	Commands       map[string]CommandWithHelp
	CommandHandler struct {
		cmds Commands
	}
)

var (
	spotifyRegex = regexp.MustCompile(`^(?:https:\/\/open.spotify.com\/track\/|spotify:track:)([a-zA-Z0-9]+)(?:.*)$`)
)

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		cmds: Commands{
			"join": CommandWithHelp{
				command: JoinCommand,
				help:    "Join a party. If no party exists in the server it will also create.",
			},
			"list": CommandWithHelp{
				command: ListCommand,
				help:    "List all users participating in the party.",
			},
			"add": CommandWithHelp{
				command: AddCommand,
				help:    "Add a new music to the users queue. You can check the musics added on Spotify.",
			},
			"sync": CommandWithHelp{
				command: SyncCommand,
				help:    "Sync music and time with the party on Spotify.",
			},
			"leave": CommandWithHelp{
				command: LeaveCommand,
				help:    "Leave the party.",
			},
			"kill": CommandWithHelp{
				command: KillCommand,
				help:    "Kill the party.",
			},
			"help": CommandWithHelp{
				command: HelpCommand,
				help:    "List of Commands.",
			},
			"info": CommandWithHelp{
				command: InfoCommand,
				help:    "Spotify Party Bot Info.",
			},
		},
	}
}

func (handler CommandHandler) Get(name string) (*Command, error) {
	cmd, found := handler.cmds[name]

	if !found {
		err := errors.New("command not found")
		return nil, err
	}

	return &cmd.command, nil
}

func JoinCommand(ctx *Context) {
	ctx.Parties.Join(ctx.Guild.ID, ctx.Channel.ID, ctx.User)

	state := ctx.Guild.ID + "-" + ctx.Channel.ID + "-" + ctx.User.ID
	url := ctx.Auth.AuthURL(state)
	dm, _ := ctx.Session.UserChannelCreate(ctx.User.ID)

	embed := &discordgo.MessageEmbed{
		Title: "Click here to connect on Spotify",
		URL:   url,
		Color: 8534465,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/wrRGQ70.png",
		},
		Description: "You need to login on Spotify to start partying with your friends! :headphones:",
	}
	ctx.Session.ChannelMessageSendEmbed(dm.ID, embed)
	message := fmt.Sprintf("<@%s> joined the party! :partying_face: ", ctx.User.ID)
	ctx.Reply(message)
}

func ListCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		askToJoinParty(ctx)
		return
	}
	fields := []*discordgo.MessageEmbedField{}
	for _, u := range party.users {
		if u.spotify == nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  ":warning: @" + u.discord.Username,
				Value: "**Not authenticated on Spotify**, please access the link sent on your DM",
			})
			continue
		}
		s, err := u.spotify.CurrentUser()
		if err != nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  ":warning: @" + u.discord.Username,
				Value: "**Not authenticated on Spotify**, please access the link sent on your DM",
			})
		} else {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  ":large_orange_diamond: @" + u.discord.Username,
				Value: "Spotify Account: @" + s.DisplayName + " ",
			})
		}
	}

	embed := discordgo.MessageEmbed{
		Title:       ":dancer: Party Guests :dancer:",
		Color:       8534465,
		Description: "Who is partying?",
		Fields:      fields,
	}
	ctx.ReplyWithEmbed(embed)
}

func AddCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		askToJoinParty(ctx)
		return
	}
	user, err := ctx.Parties.GetUser(ctx.Guild.ID, ctx.User.ID)

	if err != nil || user == nil || user.spotify == nil {
		askToJoinParty(ctx)
		return
	}

	term := strings.Join(ctx.Args, " ")
	if spotifyRegex.MatchString(term) {
		trackID := spotifyRegex.FindStringSubmatch(term)[1]
		track, err := user.spotify.GetTrack(spotify.ID(trackID))

		if err != nil {
			ctx.Reply("Track not found :cold_sweat: ")
			return
		}
		party.Add(user, *track)
		ctx.Reply("Added " + track.ExternalURLs["spotify"])
		return
	}

	results, err := user.spotify.Search(term, spotify.SearchTypeTrack)
	if err != nil || results.Tracks == nil || results.Tracks.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		ctx.Reply("track not found")
		return
	}

	track := results.Tracks.Tracks[0]
	party.Add(user, track)
	ctx.Reply("Added " + track.ExternalURLs["spotify"])
}

func SyncCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		askToJoinParty(ctx)
		return
	}
	user, err := ctx.Parties.GetUser(ctx.Guild.ID, ctx.User.ID)
	if err != nil || user == nil || user.spotify == nil {
		askToJoinParty(ctx)
		return
	}
	party.Sync(user)
	ctx.Reply(":play_pause: synced")
}

func LeaveCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		askToJoinParty(ctx)
		return
	}
	user, err := ctx.Parties.GetUser(ctx.Guild.ID, ctx.User.ID)
	if err != nil || user == nil || user.spotify == nil {
		askToJoinParty(ctx)
		return
	}
	party.Leave(user)
	ctx.Reply(fmt.Sprintf("Goodbye <@%s> :wave:", ctx.User.ID))
}

func KillCommand(ctx *Context) {
	ctx.Parties.Kill(ctx.Guild.ID, ctx.Channel.ID)
	ctx.Reply("The party is over :wave: ")
}

func HelpCommand(ctx *Context) {
	fields := []*discordgo.MessageEmbedField{}

	cmds := CmdHandler.cmds

	for cmdName, cmd := range cmds {
		fields = append(fields, &discordgo.MessageEmbedField{Name: PREFIX + cmdName, Value: cmd.help})
	}

	embed := discordgo.MessageEmbed{
		Title:       "List of Commands",
		Description: "For more info access link",
		Color:       8534465,
		Fields:      fields,
	}

	ctx.ReplyWithEmbed(embed)
}

func InfoCommand(ctx *Context) {
	embed := discordgo.MessageEmbed{
		Title:       "Spotify Party Bot",
		Description: "[Info](https://github.com/thayannevls/SpotifyPartyBot)\n[Invite](https://discord.com/api/oauth2/authorize?client_id=747675945644851221&permissions=19520&scope=bot)\n[Donate](https://ko-fi.com/thayannevls)",
		Color:       8534465,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/wrRGQ70.png",
		},
	}
	ctx.ReplyWithEmbed(embed)

}

func askToJoinParty(ctx *Context) {
	embed := discordgo.MessageEmbed{
		Title:       "Join the party first!",
		Description: fmt.Sprintf("Hey stranger :detective: <@%s>, you need to fully join a party to perform this action.", ctx.User.ID),
		Color:       8534465,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "1. Join the party",
				Value: "Type sp.join to join a current party or create a new one",
			},
			&discordgo.MessageEmbedField{
				Name:  "2. Login on Spotify",
				Value: "A link will be sent on your DM to login on Spotify.",
			},
			&discordgo.MessageEmbedField{
				Name:  "3. Have fun! :man_dancing: ",
				Value: "Have fun listening songs with your friends.",
			},
		},
	}

	ctx.ReplyWithEmbed(embed)
}
