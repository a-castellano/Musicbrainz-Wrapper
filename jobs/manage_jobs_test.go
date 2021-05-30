// +build integration_tests unit_tests

package jobs

import (
	"bytes"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestProcessJobEmptyData(t *testing.T) {

	var emptyData []byte

	origin := "MisicBrainzWrapper"

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:29:36.962Z","count":0,"offset":0,"artists":[]}
	`))}}}

	die, jobResult, err := ProcessJob(emptyData, origin, client)

	if err.Error() != "Empty job data received." {
		t.Errorf("Message with failed data should return 'Empty data received.' error, not '%s'.", err.Error())
	}

	if die == true {
		t.Errorf("Message with failed data does not stop this service.")
	}

	if len(jobResult) == 0 {
		t.Errorf("jobResult should be empty")
	}
}

func TestProcessJobErrorOnArtist(t *testing.T) {

	var infoRetrieval commontypes.InfoRetrieval
	var job commontypes.Job

	infoRetrieval.Type = commontypes.ArtistName
	infoRetrieval.Artist = "Burzum"

	retrievalData, _ := commontypes.EncodeInfoRetrieval(infoRetrieval)

	job.Data = retrievalData
	job.ID = 0
	job.Status = true
	job.Finished = false
	job.Type = 1 // https://musicmanager.gitpages.windmaker.net/Music-Manager-Docs/common-types/#job

	encodedJob, _ := commontypes.EncodeJob(job)

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-05-30T10:52:00.111Z","count":5,"offset":0,"artists":[{"id":"49cd96a6-42c3-44f6-ba2a-cd9301046b96","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":100,"name":"Burzum","sort-name":"Burzum","area":{"id":"6943f80d-1bd7-4495-8d05-5dc8da1fa6b2","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Limousis","sort-name":"Limousis","life-span":{"ended":null}},"begin-area":{"id":"34268a86-54ec-487d-afa3-1afe3266a382","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Bergen","sort-name":"Bergen","life-span":{"ended":null}},"disambiguation":"Norwegian black metal","life-span":{"begin":"1987","end":"2018-06-01","ended":true},"tags":[{"count":1,"name":"electronic"},{"count":3,"name":"norwegian"},{"count":1,"name":"ambient"},{"count":2,"name":"metal"},{"count":15,"name":"black metal"},{"count":3,"name":"dark ambient"},{"count":1,"name":"black ambient"},{"count":5,"name":"atmospheric black metal"},{"count":2,"name":"ambient black metal"},{"count":8,"name":"dungeon synth"},{"count":1,"name":"skaldic metal"}]},{"id":"9bd9798a-1ea1-4a6e-a209-bf716a1793c3","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":42,"name":"Burzumennuz","sort-name":"Burzumennuz","life-span":{"begin":"2004","ended":null}},{"id":"1587ba76-5e5c-42f5-9eaa-0e3a6b5f10ff","score":39,"name":"Alberich","sort-name":"Alberich","disambiguation":"black metal, Burzum tribute","life-span":{"ended":null}},{"id":"b9fd3422-6119-4602-915d-a7bd522d254a","score":39,"name":"Nargothrond","sort-name":"Nargothrond","disambiguation":"black metal, Burzum tribute","life-span":{"ended":null}},{"id":"d33d212f-359c-4fd0-9c6d-5fbf88e5c9dd","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":34,"name":"Uruk-Hai","sort-name":"Uruk-Hai","area":{"id":"34268a86-54ec-487d-afa3-1afe3266a382","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Bergen","sort-name":"Bergen","life-span":{"ended":null}},"begin-area":{"id":"34268a86-54ec-487d-afa3-1afe3266a382","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Bergen","sort-name":"Bergen","life-span":{"ended":null}},"end-area":{"id":"34268a86-54ec-487d-afa3-1afe3266a382","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Bergen","sort-name":"Bergen","life-span":{"ended":null}},"disambiguation":"Norwegian black metal band, precursor to Burzum","life-span":{"begin":"1988","end":"1990","ended":true}}]
	`))}}}

	origin := "MisicBrainzWrapper"
	die, jobResult, err := ProcessJob(encodedJob, origin, client)

	if err != nil {
		if !strings.HasPrefix(err.Error(), "Artist retrieval failed: unexpected end of JSON input") {
			t.Errorf("Message with failed data should return 'Artist retrieval failed: unexpected end of JSON input' error, not '%s'.", err.Error())
		}
	}

	if die == true {
		t.Errorf("Message with failed data does not stop this service.")
	}

	if len(jobResult) == 0 {
		t.Errorf("jobResult shouldn't be empty")
	}

	if len(jobResult) == 0 {
		t.Errorf("jobResult shouldn't be empty")
	}
	decodedJob, _ := commontypes.DecodeJob(jobResult)

	if decodedJob.Error != "Artist retrieval failed: unexpected end of JSON input" {
		t.Errorf("decodedJob.Error should be 'Artist retrieval failed: unexpected end of JSON input', not '%s'.", decodedJob.Error)
	}

	if decodedJob.LastOrigin != origin {
		t.Errorf("decodedJob.LastOrigin should be '%s', not '%s'.", origin, decodedJob.LastOrigin)
	}

}
