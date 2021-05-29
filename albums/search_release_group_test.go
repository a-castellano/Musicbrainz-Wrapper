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

	if err.Error() != "No release group was found." {
		t.Errorf("Error should be 'No release group was found.', not '%s'.", err)
	}

}

func TestSearchReleaseGroupWithResultsButNoMatches(t *testing.T) {

	searchAlbumInfo := MockOneReleaseGroup{}
	releaseGroup, releaseGroups, err := getReleaseGroup(searchAlbumInfo, "AnyNonExistAlbum", "AnyNonExistAlbum")

	if err != nil {
		t.Errorf("TestSearchReleaseGroupWithResultsButNoMatches should not fail.")
	}

	if releaseGroup.ID != "" {
		t.Errorf("First ReleaseGroup should have no data, not %s.", releaseGroup.ID)
	}

	if len(releaseGroups) != 0 {
		t.Errorf("releaseGroups should be an empty array.")
	}

}

func TestSearchReleaseGroupWithOnlyOneResult(t *testing.T) {

	searchAlbumInfo := MockOneReleaseGroup{}
	releaseGroup, releaseGroups, err := getReleaseGroup(searchAlbumInfo, "Hexndeifl", "Hexndeifl")

	if err != nil {
		t.Errorf("TestSearchReleaseGroupWithOnlyOneResult should not fail.")
	}

	if releaseGroup.ID != "4a9421e6-9def-4616-82bf-674d5a5c7a29" {
		t.Errorf("First ReleaseGroup ID should be '4a9421e6-9def-4616-82bf-674d5a5c7a29', not '%s'.", releaseGroup.ID)
	}

	if releaseGroup.Title != "Hexndeifl" {
		t.Errorf("First ReleaseGroup Title should be 'Hexndeifl', not '%s'.", releaseGroup.Title)
	}

	if releaseGroup.ReleaseYear != 2021 {
		t.Errorf("First ReleaseGroup ReleaseYear should be 2021, not '%d'.", releaseGroup.ReleaseYear)
	}

	if len(releaseGroups) != 0 {
		t.Errorf("releaseGroups should be an empty array.")
	}

}

func TestSearchReleaseGroupWithMoreThanOneResult(t *testing.T) {

	searchAlbumInfo := MockTwoReleaseGroups{}
	releaseGroup, releaseGroups, err := getReleaseGroup(searchAlbumInfo, "In Times Before the Light", "In Times Before the Light")

	if err != nil {
		t.Errorf("TestSearchReleaseGroupWithMoreThanOneResult should not fail.")
	}

	if releaseGroup.ID != "ba62217c-c12f-4eae-859e-5202936914ff" {
		t.Errorf("First ReleaseGroup ID should be 'ba62217c-c12f-4eae-859e-5202936914ff', not '%s'.", releaseGroup.ID)
	}

	if releaseGroup.Title != "In Times Before the Light" {
		t.Errorf("First ReleaseGroup Title should be 'In Times Before the Light', not '%s'.", releaseGroup.Title)
	}

	if releaseGroup.ReleaseYear != 2002 {
		t.Errorf("First ReleaseGroup ReleaseYear should be 2002, not '%d'.", releaseGroup.ReleaseYear)
	}

	if len(releaseGroups) != 1 {
		t.Errorf("releaseGroups should be an array with only one element.")
	}

	secondRelease := releaseGroups[0]

	if secondRelease.ID != "422496f0-d2b5-30a6-806e-cf5b7e70c28f" {
		t.Errorf("First ReleaseGroup ID should be '422496f0-d2b5-30a6-806e-cf5b7e70c28f', not '%s'.", secondRelease.ID)
	}

	if secondRelease.Title != "In Times Before the Light" {
		t.Errorf("First ReleaseGroup Title should be 'In Times Before the Light', not '%s'.", secondRelease.Title)
	}

	if secondRelease.ReleaseYear != 1997 {
		t.Errorf("First ReleaseGroup ReleaseYear should be 1997, not '%d'.", secondRelease.ReleaseYear)
	}
}
