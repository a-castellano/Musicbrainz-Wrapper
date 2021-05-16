package main

import (
	"fmt"
	albums "github.com/a-castellano/music-manager-musicbrainz-wrapper/albums"
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

	searchAlbumInfo := albums.SearchAlbumInfo{Client: client}
	//	release, _, _ := albums.SearchAlbum(searchAlbumInfo, "20:11")
	release, _, _ := albums.SearchAlbum(searchAlbumInfo, "The call")
	fmt.Println(release)
}
