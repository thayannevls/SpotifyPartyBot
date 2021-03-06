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
	guildID, userID := s[0], s[2]
	user, err := Parties.GetUser(guildID, userID)
	if err != nil {
		http.Error(w, "Error trying to retrieve user", http.StatusNotFound)
		fmt.Println("Error trying to retrieve user: ", err)
		return
	}
	userSpotify := Auth.NewClient(token)
	updatedUserWithSpotify := NewUser(user.discord, &userSpotify)
	updatedUserWithSpotify.CreatePlaylist()
	party := Parties.GetByGuild(guildID)
	Parties.UpdateUser(party, user, updatedUserWithSpotify)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully authenticated. Ready to party!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Listening on "+port)
}

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

		log.Println("Signal type : ", signalType)

		os.Exit(0)

	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, r)

	wg.Done()
}
