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

func readSearchedReleaseGroup(body []byte, album string) (ReleaseGroup, []ReleaseGroup, error) {

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
		return mainResult, otherResults, errors.New("No release group was found.")
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

func getReleaseGroup(searchAlbumInfo SearchAlbumInfoInterface, album string, albumString string) (ReleaseGroup, []ReleaseGroup, error) {

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

	releaseGroup, extraReleaseGroups, releaseGroupErr := readSearchedReleaseGroup(releaseGroupBody, album)
	if releaseGroupErr != nil {
		return releaseGroup, extraReleaseGroups, releaseGroupErr
	}

	return releaseGroup, extraReleaseGroups, nil
}
