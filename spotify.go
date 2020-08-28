package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
)

// Search ..
func Search(term string) (track spotify.FullTrack, err error) {
	results, err := client.Search(term, spotify.SearchTypeTrack)

	if err != nil {
		log.Fatal(err)
	}

	if results.Tracks != nil {
		if results.Tracks.Tracks == nil {
			fmt.Println(" ---  ", results)
			return
		}
		for _, item := range results.Tracks.Tracks {
			fmt.Println("   ", item.Name)
		}
		return results.Tracks.Tracks[0], nil
	}

	err = errors.New("couldn't found track")
	return
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
