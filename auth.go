package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["state"]
	if !ok || len(params[0]) < 1 {
		log.Println("Url Param 'state' is missing")
		return
	}
	state := params[0]

	token, err := Auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	s := strings.Split(state, "-")
	guildID, channelID, userID := s[0], s[1], s[2]
	user, err := Parties.GetUser(guildID, channelID, userID)
	if err != nil {
		http.Error(w, "Error trying to retrieve user", http.StatusNotFound)
		fmt.Println("Error trying to retrieve user: ", err)
		return
	}
	userSpotify := Auth.NewClient(token)
	updatedUserWithSpotify := NewUser(user.discord, &userSpotify)
	party := Parties.GetByGuild(guildID)
	Parties.UpdateUser(party, user, updatedUserWithSpotify)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully authenticated. Ready to party!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Listening....")
}

// InitAuthServer ...
func InitAuthServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/auth", authHandler)
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		signalType := <-ch
		signal.Stop(ch)
		log.Println("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		log.Println("Signal type : ", signalType)

		os.Exit(0)

	}()

	http.ListenAndServe(":8080", r)

	wg.Done()
}
