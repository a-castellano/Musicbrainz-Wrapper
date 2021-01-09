package main

import (
	artists "github.com/a-castellano/music-manager-musicbrainz-wrapper/artists"
	"log"
	"net/http"
	"time"
)

func main() {

	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	log.Println("Dummy test")

	artists.SearchArtist(client, "Melechesh")

}
