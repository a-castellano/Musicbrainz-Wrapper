package albums

import (
	"net/http"
)

type Track struct {
	ID       string
	Title    string
	Lenght   int
	Position int
}

type ReleaseGroup struct {
	ID          string
	Title       string
	ReleaseYear int
	Releases    []Release
}

type Release struct {
	ID     string
	Title  string
	Tracks []Track
}

type SearchAlbumInfoInterface interface {
	SearchReleaseGroups(req *http.Request) (*http.Response, error)
	GetReleases(req *http.Request) (*http.Response, error)
	GetReleaseInfo(req *http.Request) (*http.Response, error)
}

type SearchAlbumInfo struct {
	Client http.Client
}

func (s SearchAlbumInfo) SearchReleaseGroups(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func (s SearchAlbumInfo) GetReleases(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func (s SearchAlbumInfo) GetReleaseInfo(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func SearchReleaseGroups(s SearchAlbumInfoInterface, req *http.Request) (*http.Response, error) {
	response, responseError := s.SearchReleaseGroups(req)
	return response, responseError
}

func GetReleases(s SearchAlbumInfoInterface, req *http.Request) (*http.Response, error) {
	response, responseError := s.GetReleases(req)
	return response, responseError
}
