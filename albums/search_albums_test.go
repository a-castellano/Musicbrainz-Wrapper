// +build integration_tests unit_tests

package albums

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSearchReleaseGroupWithNoResults(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-17T11:29:49.260Z","count":0,"offset":0,"release-groups":[]}
	`))}}}

	searchAlbumInfo := SearchAlbumInfo{Client: client}
	_, _, err := SearchAlbum(searchAlbumInfo, "AnyNonExistAlbum")

	if err == nil {
		t.Errorf("TestSearchReleaseGroupWithNoResult should fail.")
	}

	//	if err.Error() != "No ReleaseGroup was found." {
	//		t.Errorf("Error should be 'No ReleaseGroup was found.', not '%s'.", err)
	//	}

}
