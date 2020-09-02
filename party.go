package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/zmb3/spotify"
)

// Party ...
type Party struct {
	guildID, channelID string
	users              map[string]*User
	queue              []*User
}

// PartyManager ...
type PartyManager struct {
	parties map[string]*Party
}

// NewPartyManager ...
func NewPartyManager() *PartyManager {
	return &PartyManager{make(map[string]*Party)}
}

// NewParty create a new party
func NewParty(guildID string, channelID string) *Party {
	party := new(Party)
	party.guildID = guildID
	party.channelID = channelID
	party.users = make(map[string]*User)
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

// Join a party in the guild
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

	return party.users[oldUser.discord.ID], nil
}

func (party *Party) Play() {
	for _, user := range party.users {
		go func(user *User) {
			user.spotify.Play()
		}(user)
	}
}

func (party *Party) Pause() {
	for _, user := range party.users {
		go func(user *User) {
			user.spotify.Pause()
		}(user)
	}
}

func (party *Party) Add(track spotify.FullTrack) {
	for _, user := range party.users {
		go func(user *User) {
			user.spotify.QueueSong(track.ID)
		}(user)
	}
}
