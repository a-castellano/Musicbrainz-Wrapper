// +build integration_tests unit_tests

package jobs

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

	origin := "MsicBrainzWrapper"

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
