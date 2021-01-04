package artists

import (
	"encoding/json"
	"errors"
	"fmt"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SearchArtistData commontypes.Artist

func processResult(searchResult map[string]interface{}) (SearchArtistData, []SearchArtistData, error) {
	var artistData SearchArtistData
	var artistExtraData []SearchArtistData

	numberOfResultsString := fmt.Sprintf("%v", searchResult["count"])
	numberOfResults, _ := strconv.Atoi(numberOfResultsString)
	if numberOfResults == 0 {
		return artistData, artistExtraData, errors.New("No artist was found.")
	}

	return artistData, artistExtraData, nil
}

func SearchArtist(client http.Client, artist string) (SearchArtistData, []SearchArtistData, error) {

	var artistData SearchArtistData
	var artistExtraData []SearchArtistData

	var searchResult map[string]interface{}

	artistString := strings.Replace(artist, " ", "%20", -1)
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/?query=%s&fmt=json", artistString)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return artistData, artistExtraData, err
	}

	req.Header.Set("User-Agent", "https://github.com/a-castellano/music-manager-musicbrainz-wrapper")

	res, getErr := client.Do(req)
	if getErr != nil {
		return artistData, artistExtraData, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return artistData, artistExtraData, readErr
	}

	jsonErr := json.Unmarshal([]byte(body), &searchResult)
	if jsonErr != nil {
		return artistData, artistExtraData, jsonErr
	}

	return processResult(searchResult)
}
