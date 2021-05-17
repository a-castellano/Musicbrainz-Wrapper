package albums

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

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
