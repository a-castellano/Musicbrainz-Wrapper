// +build integration_tests unit_tests

package albums

import (
	"testing"
)

func TestGetReleasesWithOneResult(t *testing.T) {

	releaseGroup := ReleaseGroup{ID: "4a9421e6-9def-4616-82bf-674d5a5c7a29", Title: "Hexndeifl", ReleaseYear: 2021}
	mockOneReleaseGroup := MockOneReleaseGroup{}

	releases, getReleasesError := getReleasesFromReleaseGroup(mockOneReleaseGroup, releaseGroup)

	if getReleasesError != nil {
		t.Errorf("TestSearchReleaseGroupWithNoResults should not fail.")
	}

	if len(releases) != 1 {
		t.Errorf("TestSearchReleaseGroupWithNoResults should return one release only.")
	}

	firstRelease := releases[0]

	if firstRelease.ID != "cb36e0ab-6634-4d85-a84e-8be89924811f" {
		t.Errorf("TestSearchReleaseGroupWithNoResults first release id should be 'cb36e0ab-6634-4d85-a84e-8be89924811f', not '%s'.", firstRelease.ID)
	}

	if firstRelease.Title != "Hexndeifl" {
		t.Errorf("TestSearchReleaseGroupWithNoResults first release title should be 'Hexndeifl', not '%s'.", firstRelease.Title)
	}

	if len(firstRelease.Tracks) != 6 {
		t.Errorf("TestSearchReleaseGroupWithNoResults first release trackNumber should be 6, not '%d'.", len(firstRelease.Tracks))
	}

}
