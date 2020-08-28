package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

var (
	playlist spotify.FullPlaylist
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	if m.Content == "!sp-join" {
		url := auth.AuthURL(state)
		dm, _ := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(dm.ID, url)
	}

	if m.Content == "!sp-start" {
		user, _ := client.CurrentUser()
		p, _ := client.CreatePlaylistForUser(user.ID, "Spotify Party :: Queue", "", true)
		playlist = *p
	}

	if m.Content == "!sp-play" {
		s.ChannelMessageSend(m.ChannelID, "Play!")
		tracks, _ := client.GetPlaylistTracks(playlist.ID)
		client.PlayOpt(&spotify.PlayOptions{PlaybackContext: &playlist.URI})
		t := tracks.Tracks[0]

		d := t.Track.TimeDuration()
		fmt.Println(d)
	}

	if m.Content == "!sp-pause" {
		client.Pause()
	}

	if m.Content == "!sp-next" {
		client.Next()
	}

	if strings.HasPrefix(m.Content, "!sp-add") {
		track, err := Search(strings.TrimPrefix(m.Content, "!sp-add"))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Couldn't find track!")
		}
		s.ChannelMessageSend(m.ChannelID, "Adding to playlist..."+track.Name)
		client.AddTracksToPlaylist(playlist.ID, track.ID)
	}

	if strings.HasPrefix(m.Content, "!sp-search") {
		track, err := Search(strings.TrimPrefix(m.Content, "!sp-search"))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Couldn't find track!")
		}

		client.PlayOpt(&spotify.PlayOptions{URIs: []spotify.URI{track.URI}})
		s.ChannelMessageSend(m.ChannelID, "Playing "+track.ExternalURLs["spotify"])

	}
}

// ConfigDiscord ...
func ConfigDiscord() {
	token := os.Getenv("DISCORD_SECRET")
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	wg.Done()
}
