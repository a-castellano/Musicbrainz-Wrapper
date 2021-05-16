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

func readSearchAlbum(body []byte, album string) (ReleaseGroup, []ReleaseGroup, error) {

	var mainResult ReleaseGroup
	var otherResults []ReleaseGroup

	var results map[string]interface{}

	albumString := strings.ToLower(album)

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
			candidateTitle := strings.ToLower(candidate["title"].(string))
			if albumString == candidateTitle {
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

	reqReleaseGroupResponse, reqReleaseGroupResponseError := searchAlbumInfo.SearchReleaseGroups(reqReleaseGroup)

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

func ReadReleaseTracks(releaseInfoBody []byte) ([]Track, error) {

	var tracks []Track
	var tracksResult map[string]interface{}

	jsonErr := json.Unmarshal([]byte(releaseInfoBody), &tracksResult)
	if jsonErr != nil {
		return tracks, jsonErr
	}

	mediaSlice := reflect.ValueOf(tracksResult["media"])
	if mediaSlice.IsValid() {
		for i := 0; i < mediaSlice.Len(); i++ {
			mediaInfo := mediaSlice.Index(i).Interface().(map[string]interface{})
			tracksByMediaInfo := reflect.ValueOf(mediaInfo["tracks"])
			for j := 0; j < tracksByMediaInfo.Len(); j++ {
				trackInfo := tracksByMediaInfo.Index(j).Interface().(map[string]interface{})
				var track Track
				track.ID = reflect.ValueOf(trackInfo["id"]).String()
				track.Title = reflect.ValueOf(trackInfo["title"]).String()
				track.Lenght = int(trackInfo["length"].(float64))
				tracks = append(tracks, track)
			}
		}
	}

	return tracks, nil
}

func getReleasesFromReleaseGroup(searchAlbumInfo SearchAlbumInfo, releaseGroup ReleaseGroup) ([]Release, error) {

	var releases []Release
	var releaseResult map[string]interface{}

	searchReleaseGroupUrl := fmt.Sprintf("https://musicbrainz.org/ws/2/release-group/%s?fmt=json&inc=releases", releaseGroup.ID)

	releaseRequest, releaseRequesterr := http.NewRequest(http.MethodGet, searchReleaseGroupUrl, nil)
	if releaseRequesterr != nil {
		return releases, releaseRequesterr
	}

	releaseRequest.Header.Set("User-Agent", "https://github.com/a-castellano/music-manager-musicbrainz-wrapper")

	reqReleaseResponse, reqReleaseResponseError := GetReleases(searchAlbumInfo, releaseRequest)

	if reqReleaseResponseError != nil {
		return releases, reqReleaseResponseError
	}

	releaseBody, releaseReadErr := ioutil.ReadAll(reqReleaseResponse.Body)
	if releaseReadErr != nil {
		return releases, releaseReadErr
	}

	jsonErr := json.Unmarshal([]byte(releaseBody), &releaseResult)
	if jsonErr != nil {
		return releases, jsonErr
	}

	releasesSlice := reflect.ValueOf(releaseResult["releases"])
	if releasesSlice.IsValid() {

		for i := 0; i < releasesSlice.Len(); i++ {
			releaseInfo := releasesSlice.Index(i).Interface().(map[string]interface{})
			releaseStatus := reflect.ValueOf(releaseInfo["status"]).String()
			if releaseStatus == "Official" {
				var release Release
				release.ID = reflect.ValueOf(releaseInfo["id"]).String()
				release.Title = reflect.ValueOf(releaseInfo["title"]).String()
				releaseQueryString := fmt.Sprintf("https://musicbrainz.org/ws/2/release/%s?fmt=json&inc=recordings", release.ID)

				reqRelease, errReqRelease := http.NewRequest(http.MethodGet, releaseQueryString, nil)
				if errReqRelease != nil {
					return releases, errReqRelease
				}

				reqRelease.Header.Set("User-Agent", "https://github.com/a-castellano/music-manager-musicbrainz-wrapper")

				releaseInfoRaw, releaseInfoErr := searchAlbumInfo.GetReleaseInfo(reqRelease)

				if releaseInfoErr != nil {
					return releases, releaseInfoErr
				}

				releaseInfoBody, releaseInfoReadErr := ioutil.ReadAll(releaseInfoRaw.Body)
				if releaseInfoReadErr != nil {
					return releases, releaseInfoReadErr
				}

				releaseTracks, releaseTracksErr := ReadReleaseTracks(releaseInfoBody)

				if releaseTracksErr != nil {
					return releases, releaseTracksErr
				}

				release.Tracks = releaseTracks

				releases = append(releases, release)

			}
		}
	}
	return releases, nil
}

func SearchAlbum(searchAlbumInfo SearchAlbumInfo, album string) (Release, []Release, error) {

	// First query for relelase groups matching album string
	//	var releaseGroup ReleaseGroup
	//	var extraReleaseGroup []ReleaseGroup
	var releases []Release
	var release Release
	var extraReleases []Release

	albumString := strings.Replace(album, " ", "%20", -1)

	releaseGroup, otherReleaseGroups, releaseGrouperr := getReleaseGroup(searchAlbumInfo, album, albumString)

	if releaseGrouperr != nil {
		return release, extraReleases, releaseGrouperr
	}

	releasesFromReleaseGroup, getReleasesErr := getReleasesFromReleaseGroup(searchAlbumInfo, releaseGroup)
	releases = releasesFromReleaseGroup

	if getReleasesErr != nil {
		return release, extraReleases, releaseGrouperr
	}

	for _, releaseGroup := range otherReleaseGroups {
		releasesFromReleaseGroup, getReleasesErr = getReleasesFromReleaseGroup(searchAlbumInfo, releaseGroup)
		if getReleasesErr != nil {
			return release, extraReleases, releaseGrouperr
		}
		releases = append(releases, releasesFromReleaseGroup...)
	}

	release = releases[0]
	extraReleases = releases[1:]

	return release, extraReleases, nil
}
