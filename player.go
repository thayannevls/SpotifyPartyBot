package main

import (
	"time"

	"github.com/zmb3/spotify"
)

type Player struct {
	running      bool
	started      time.Time
	duration     time.Duration
	currentTrack *spotify.FullTrack
}

func NewPlayer() *Player {
	player := new(Player)
	player.running = false
	return player
}

func (player *Player) Play(track *spotify.FullTrack, callback func()) {
	player.running = true
	player.duration = track.TimeDuration()
	player.started = time.Now()
	player.currentTrack = track

	time.Sleep(player.duration)

	player.running = false
	player.started = *new(time.Time)
	player.duration = *new(time.Duration)
	player.currentTrack = nil

	callback()
}
