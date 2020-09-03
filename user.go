package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

type User struct {
	discord  *discordgo.User
	spotify  *spotify.Client
	playlist *spotify.FullPlaylist
}

func NewUser(discord *discordgo.User, spotify *spotify.Client) *User {
	user := new(User)
	user.discord = discord
	user.spotify = spotify

	return user
}

func (user *User) CreatePlaylist() {
	u, err := user.spotify.CurrentUser()
	if err != nil {
		fmt.Println("Error getting user")
	}
	p, err := user.spotify.CreatePlaylistForUser(u.ID, "Spotify Party :: Queue", "", true)
	if err != nil {
		fmt.Println("Error creating playlist")
	}
	user.playlist = p
}

func (user *User) PopFromPlaylist() (*spotify.FullTrack, error) {
	tracks, err := user.spotify.GetPlaylistTracks(user.playlist.ID)
	if err != nil || tracks == nil || tracks.Tracks == nil || len(tracks.Tracks) == 0 {
		err := errors.New("something happened when getting the playlist")
		return nil, err
	}

	track := tracks.Tracks[0].Track
	user.spotify.RemoveTracksFromPlaylist(user.playlist.ID, track.ID)

	return &track, nil
}
