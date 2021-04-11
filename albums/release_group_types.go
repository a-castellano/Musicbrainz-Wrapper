package albums

import (
	"net/http"
)

type ReleaseGroup struct {
	ID          string
	Title       string
	ReleaseYear int
	Releases    []Release
}

type Release struct {
	ID    string
	Title string
}

func SearchReleaseGroup(client http.Client, album string) (ReleaseGroup, []ReleaseGroup, error) {
	var releaseGroup ReleaseGroup
	var extraReleaseGroup []ReleaseGroup

	//https: //musicbrainz.org/ws/2/release-group/?query=confrontaskdas%C3%B1jdhasjkldha&fmt=json

	return releaseGroup, extraReleaseGroup, nil
}
