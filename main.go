package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

var (
	state     = "abc1234"
	auth      spotify.Authenticator
	client    spotify.Client
	wg        = new(sync.WaitGroup)
	botPrefix = "!sp-"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wg.Add(2)
	auth = GetSpotifyAuth()
	go ConfigDiscord()
	go InitAuthServer()

	wg.Wait()
}
