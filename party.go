package main

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

type Party struct {
	guildID, channelID string
	users              map[string]*User
	current            int
	queue              []*User
	player             *Player
}

type PartyManager struct {
	parties map[string]*Party
}

func NewPartyManager() *PartyManager {
	return &PartyManager{make(map[string]*Party)}
}

func NewParty(guildID string, channelID string) *Party {
	party := new(Party)
	party.guildID = guildID
	party.channelID = channelID
	party.users = make(map[string]*User)
	party.queue = []*User{}
	party.current = 0
	party.player = NewPlayer()
	return party
}

// GetByGuild get party by guildID
func (manager *PartyManager) GetByGuild(guildID string) *Party {
	for _, party := range manager.parties {
		if party.guildID == guildID {
			return party
		}
	}
	return nil
}

func (manager *PartyManager) Join(guildID, channelID string, userDiscord *discordgo.User) *Party {
	user := NewUser(userDiscord, nil)
	party := manager.GetByGuild(guildID)
	if party == nil {
		party = NewParty(guildID, channelID)
		manager.parties[channelID] = party
	}

	party.users[userDiscord.ID] = user
	return party
}

func (manager *PartyManager) Kill(guildID, channelID string) {
	party := manager.GetByGuild(guildID)
	if party == nil {
		return
	}
	party.Pause()
	party = nil
	delete(manager.parties, channelID)
}

func (manager *PartyManager) GetUser(guildID, userID string) (*User, error) {
	party := manager.GetByGuild(guildID)

	if party == nil {
		err := errors.New("party not found")
		return nil, err
	}

	if party.users[userID] == nil {
		err := errors.New("user not found in party")
		return nil, err
	}

	return party.users[userID], nil
}

func (manager *PartyManager) UpdateUser(party *Party, oldUser, newUser *User) (*User, error) {
	if party == nil || oldUser == nil {
		err := errors.New("received invalid parameter")
		return nil, err
	}
	if oldUser.discord.ID != newUser.discord.ID {
		err := errors.New("oldUser and newUser IDs did not match")
		return nil, err
	}
	if party.users[oldUser.discord.ID] == nil {
		err := errors.New("user not found in party")
		return nil, err
	}

	party.users[oldUser.discord.ID] = newUser
	party.queue = append(party.queue, newUser)

	return party.users[oldUser.discord.ID], nil
}

func (party *Party) Play() {
	party.PlayAux(0)
}
func (party *Party) PlayAux(notFound int) {
	if len(party.queue) == 0 || notFound == len(party.queue) {
		party.Pause()
		return
	}
	u := party.queue[party.current]

	track, err := u.PopFromPlaylist()

	if err != nil {
		party.current = (party.current + 1) % len(party.queue)
		party.PlayAux(notFound + 1)
		return
	}

	for _, user := range party.users {
		if user == nil || user.spotify == nil {
			continue
		}
		go func(user *User) {
			user.spotify.PlayOpt(&spotify.PlayOptions{URIs: []spotify.URI{track.URI}})
		}(user)
	}

	go party.player.Play(track, func() {
		if len(party.queue) == 0 {
			party.Pause()
			return
		}
		party.current = (party.current + 1) % len(party.queue)
		party.PlayAux(0)
	})
}

func (party *Party) Stop() {
	party.player.running = false
	party.Pause()
}

func (party *Party) Pause() {
	for _, user := range party.users {
		if user == nil || user.spotify == nil {
			continue
		}
		go func(user *User) {
			user.spotify.Pause()
		}(user)
	}
}

func (party *Party) Add(user *User, track spotify.FullTrack) {
	user.spotify.AddTracksToPlaylist(user.playlist.ID, track.ID)

	if party.player.running {
		return
	}
	party.Play()
}

func (party *Party) Sync(user *User) {
	currentTrack := party.player.currentTrack
	duration := party.player.duration
	started := party.player.started
	positionMs := diffDate(started, duration)
	user.spotify.PlayOpt(&spotify.PlayOptions{URIs: []spotify.URI{currentTrack.URI}, PositionMs: positionMs})
}

func (party *Party) Leave(user *User) {
	index := -1
	for i, u := range party.queue {
		if u.discord.ID == user.discord.ID {
			index = i
			break
		}
	}
	if index != -1 {
		party.queue = append(party.queue[:index], party.queue[index+1:]...)
	}
	delete(party.users, user.discord.ID)
}

func diffDate(start time.Time, duration time.Duration) int {
	now := time.Now()

	sub := now.Sub(start)

	return int(sub / time.Millisecond)
}
