// +build integration_tests unit_tests

package artists

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSearchArtistWithNoResults(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:29:36.962Z","count":0,"offset":0,"artists":[]}
	`))}}}

	_, _, err := SearchArtist(client, "AnyNonExistentArtist")

	if err == nil {
		t.Errorf("TestClientNoArtists should fail.")
	}

	if err.Error() != "No artist was found." {
		t.Errorf("Error should be 'No artist was found.', not '%s'.", err)
	}

}
