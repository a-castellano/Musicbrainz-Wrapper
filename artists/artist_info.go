package artists

import (
	"encoding/json"
	"fmt"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"io/ioutil"
	"net/http"
	"reflect"
)

func obtainRecordInfo(info interface{}) commontypes.Record {
	var record commontypes.Record

	return record
}

func GetArtistRecords(client http.Client, artistData SearchArtistData) ([]commontypes.Record, error) {
	var records []commontypes.Record
	var artistInfo map[string]interface{}

	url := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/%s?fmt=json&inc=release-groups", artistData.ID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return records, err
	}

	req.Header.Set("User-Agent", "https://github.com/a-castellano/music-manager-musicbrainz-wrapper")

	res, getErr := client.Do(req)
	if getErr != nil {
		return records, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return records, readErr
	}

	jsonErr := json.Unmarshal([]byte(body), &artistInfo)
	if jsonErr != nil {
		return records, jsonErr
	}

	releaseGroups := reflect.ValueOf(artistInfo["release-groups"])
	for i := 0; i < releaseGroups.Len(); i++ {
		releaseInfo := releaseGroups.Index(i).Interface().(map[string]interface{})
		record := obtainRecordInfo(releaseInfo)
		records = append(records, record)
		fmt.Println(releaseInfo)
	}

	return records, nil
}
