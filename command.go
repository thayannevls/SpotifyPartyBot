package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zmb3/spotify"
)

type (
	Command        func(*Context)
	Commands       map[string]Command
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
			"join":  JoinCommand,
			"list":  ListCommand,
			"play":  PlayCommand,
			"pause": PauseCommand,
			"add":   AddCommand,
			"start": StartCommand,
			"sync":  SyncCommand,
		},
	}
}

func (handler CommandHandler) Get(name string) (*Command, error) {
	cmd, found := handler.cmds[name]

	if !found {
		err := errors.New("command not found")
		return nil, err
	}

	return &cmd, nil
}

func JoinCommand(ctx *Context) {
	ctx.Parties.Join(ctx.Guild.ID, ctx.Channel.ID, ctx.User)

	state := ctx.Guild.ID + "-" + ctx.Channel.ID + "-" + ctx.User.ID
	url := ctx.Auth.AuthURL(state)
	dm, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
	ctx.Session.ChannelMessageSend(dm.ID, url)
}

func ListCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	reply := ctx.Guild.Name + " party's list"
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}

	for _, u := range party.users {
		if u.spotify == nil {
			reply += "\n - " + u.discord.Username + " - Not authenticated, please join the party."
			continue
		}
		s, err := u.spotify.CurrentUser()
		if err != nil {
			reply += "\n - " + u.discord.Username + " - Not authenticated, please join the party."
		} else {
			reply += "\n - " + u.discord.Username + " / " + s.DisplayName
		}
	}
	ctx.Reply(reply)
}

func PlayCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}
	party.Play()
	ctx.Reply("Playing!")
}

func PauseCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}

	party.Pause()
	ctx.Reply("Paused!")
}

func AddCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}
	user, err := ctx.Parties.GetUser(ctx.Guild.ID, ctx.User.ID)

	if err != nil {
		ctx.Reply("join the party to add music")
		return
	}

	term := strings.Join(ctx.Args, " ")
	if spotifyRegex.MatchString(term) {
		trackID := spotifyRegex.FindStringSubmatch(term)[1]
		track, err := user.spotify.GetTrack(spotify.ID(trackID))

		if err != nil {
			ctx.Reply("track not found...: ")
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

func StartCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}

	party.Start()
	ctx.Reply("Party started!")
}

func SyncCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	if party == nil {
		ctx.Reply("no party is happening, type sp.join to start one")
		return
	}
	user, err := ctx.Parties.GetUser(ctx.Guild.ID, ctx.User.ID)
	if err != nil {
		ctx.Reply("join the party first please")
		return
	}
	party.Sync(user)
	ctx.Reply("Sync!")
}
