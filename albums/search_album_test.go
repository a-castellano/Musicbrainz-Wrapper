// +build integration_tests unit_tests

package albums

import (
	"testing"
)

func TestSearchAlbumWithOneResult(t *testing.T) {

	mockOneReleaseGroup := MockOneReleaseGroup{}

	release, releases, searchAlbumErr := SearchAlbum(mockOneReleaseGroup, "Hexndeifl")

	if searchAlbumErr != nil {
		t.Errorf("TestSearchAlbumWithOneResult shouldn't fail.")
	}

	if len(releases) != 0 {
		t.Errorf("TestSearchAlbumWithOneResult should return only one release, not %d.", len(releases))
	}

	if release.ID != "cb36e0ab-6634-4d85-a84e-8be89924811f" {
		t.Errorf("Release ID should be 'cb36e0ab-6634-4d85-a84e-8be89924811f', not %s.", release.ID)
	}

	if release.Title != "Hexndeifl" {
		t.Errorf("Release Title should be 'Hexndeifl', not %s.", release.Title)
	}

	if len(release.Tracks) != 6 {
		t.Errorf("Release should contain 6 tracks, not %d.", len(release.Tracks))
	}

	track := release.Tracks[0]

	if track.Title != "Transilvania" {
		t.Errorf("First track title should be 'Transilvania', not %s.", track.Title)
	}
}
