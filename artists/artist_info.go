package artists

import (
	"encoding/json"
	"fmt"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func calculateRecordType(recordInfo map[string]interface{}) commontypes.RecordType {
	var recordType commontypes.RecordType
	var secondaryType string = ""
	switch reflect.ValueOf(recordInfo["primary-type"]).String() {
	case "Album":
		recordType = commontypes.FullLength
	case "EP":
		recordType = commontypes.EP
	case "Single":
		recordType = commontypes.Single
	default:
		recordType = commontypes.Other
	}
	if reflect.ValueOf(recordInfo["secondary-types"]).Len() != 0 {
		secondaryType = reflect.ValueOf(recordInfo["secondary-types"]).Index(0).Interface().(string)
		switch secondaryType {
		case "Compilation":
			recordType = commontypes.Compilation
		case "Live":
			recordType = commontypes.Live
		case "Demo":
			recordType = commontypes.Demo
		}
	}

	return recordType
}

func obtainRecordInfo(info interface{}) commontypes.Record {
	var record commontypes.Record

	recordInfo := info.(map[string]interface{})

	record.Name = reflect.ValueOf(recordInfo["title"]).String()
	record.ID = reflect.ValueOf(recordInfo["id"]).String()
	record.URL = fmt.Sprintf("https://musicbrainz.org/release-group/%s", reflect.ValueOf(recordInfo["id"]).String())
	record.Year, _ = strconv.Atoi(strings.Split(reflect.ValueOf(recordInfo["first-release-date"]).String(), "-")[0])
	record.Type = calculateRecordType(recordInfo)

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
	}

	return records, nil
}
