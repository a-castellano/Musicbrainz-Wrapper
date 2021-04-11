package albums

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
	ReleaseGroups(req *http.Request) (*http.Response, error)
	Releases(req *http.Request) (*http.Response, error)
}

type SearchAlbumInfo struct {
	Client http.Client
}

func (s SearchAlbumInfo) SearchReleaseGroups(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func (s SearchAlbumInfo) ReleaseGroups(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func (s SearchAlbumInfo) Releases(req *http.Request) (*http.Response, error) {
	response, responseError := s.Client.Do(req)

	return response, responseError

}

func SearchReleaseGroups(s SearchAlbumInfoInterface, req *http.Request) (*http.Response, error) {
	response, responseError := s.SearchReleaseGroups(req)
	return response, responseError
}

func readSearchAlbum(body []byte, album string) (ReleaseGroup, []ReleaseGroup, error) {

	var mainResult ReleaseGroup
	var otherResults []ReleaseGroup

	var results map[string]interface{}

	jsonErr := json.Unmarshal([]byte(body), &results)
	if jsonErr != nil {
		return mainResult, otherResults, jsonErr
	}

	reflectedNumberOfResults := reflect.ValueOf(results["count"])
	numberOfResults := int(reflectedNumberOfResults.Interface().(float64))
	if numberOfResults == 0 {
		return mainResult, otherResults, errors.New("No album was found.")
	} else {

		releaseGroupSlice := reflect.ValueOf(results["release-groups"])
		for i := 0; i < releaseGroupSlice.Len(); i++ {
			candidate := releaseGroupSlice.Index(i).Interface().(map[string]interface{})
			if album == candidate["title"] {
				score := int(candidate["score"].(float64))
				if score == 100 && mainResult.ID == "" { // First ocurrence
					mainResult.ID = reflect.ValueOf(candidate["id"]).String()
					mainResult.Title = reflect.ValueOf(candidate["title"]).String()
					releaseDate := reflect.ValueOf(candidate["first-release-date"]).String()
					releaseYear, _ := strconv.Atoi(strings.Split(releaseDate, "-")[0])
					mainResult.ReleaseYear = releaseYear
				} else {
					if score == 100 {
						var extraReleaseGroup ReleaseGroup
						extraReleaseGroup.ID = reflect.ValueOf(candidate["id"]).String()
						extraReleaseGroup.Title = reflect.ValueOf(candidate["title"]).String()
						releaseDate := reflect.ValueOf(candidate["first-release-date"]).String()
						releaseYear, _ := strconv.Atoi(strings.Split(releaseDate, "-")[0])
						extraReleaseGroup.ReleaseYear = releaseYear
						otherResults = append(otherResults, extraReleaseGroup)
					}
				}
			}
		}

	}

	return mainResult, otherResults, nil
}

func getReleaseGroup(searchAlbumInfo SearchAlbumInfo, album string, albumString string) (ReleaseGroup, []ReleaseGroup, error) {

	var releaseGroup ReleaseGroup
	var extraReleaseGroups []ReleaseGroup

	searchReleaseGroupUrl := fmt.Sprintf("https://musicbrainz.org/ws/2/release-group/?query=%s&fmt=json", albumString)

	reqReleaseGroup, errReqReleaseGroup := http.NewRequest(http.MethodGet, searchReleaseGroupUrl, nil)
	if errReqReleaseGroup != nil {
		return releaseGroup, extraReleaseGroups, errReqReleaseGroup
	}

	reqReleaseGroup.Header.Set("User-Agent", "https://github.com/a-castellano/music-manager-musicbrainz-wrapper")

	reqReleaseGroupResponse, reqReleaseGroupResponseError := SearchReleaseGroups(searchAlbumInfo, reqReleaseGroup)

	if reqReleaseGroupResponseError != nil {
		return releaseGroup, extraReleaseGroups, reqReleaseGroupResponseError
	}

	releaseGroupBody, releaseGroupReadErr := ioutil.ReadAll(reqReleaseGroupResponse.Body)
	if releaseGroupReadErr != nil {
		return releaseGroup, extraReleaseGroups, releaseGroupReadErr
	}

	releaseGroup, extraReleaseGroups, releaseGroupErr := readSearchAlbum(releaseGroupBody, album)
	if releaseGroupErr != nil {
		return releaseGroup, extraReleaseGroups, releaseGroupErr
	}

	return releaseGroup, extraReleaseGroups, nil
}

func SearchAlbum(searchAlbumInfo SearchAlbumInfo, album string) (Release, []Release, error) {

	// First query for relelase groups matching album string
	//	var releaseGroup ReleaseGroup
	//	var extraReleaseGroup []ReleaseGroup
	var release Release
	var extraReleases []Release

	albumString := strings.Replace(album, " ", "%20", -1)

	releaseGroup, otherReleaseGroups, releaseGrouperr := getReleaseGroup(searchAlbumInfo, album, albumString)

	if releaseGrouperr != nil {
		return release, extraReleases, releaseGrouperr
	}
	fmt.Println(releaseGroup)
	fmt.Println(otherReleaseGroups)

	release = getReleaseFromReleaseGroup(searchAlbumInfo, releaseGroup)
	//	fmt.Println(releaseGroupBody)
	//https://musicbrainz.org/ws/2/release-group/495064c7-a65f-36f6-952d-c0990222d459?fmt=json&inc=releases
	//https://musicbrainz.org/ws/2/release/1b704279-f088-4df7-aed9-35c57e79ae15?fmt=json&inc=recordings

	return release, extraReleases, nil
}
