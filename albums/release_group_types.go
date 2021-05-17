package albums

import (
	"strings"
)

func SearchAlbum(searchAlbumInfo SearchAlbumInfo, album string) (Release, []Release, error) {

	var releases []Release
	var release Release
	var extraReleases []Release

	albumString := strings.Replace(album, " ", "%20", -1)

	releaseGroup, otherReleaseGroups, releaseGrouperr := getReleaseGroup(searchAlbumInfo, album, albumString)

	if releaseGrouperr != nil {
		return release, extraReleases, releaseGrouperr
	}

	releasesFromReleaseGroup, getReleasesErr := getReleasesFromReleaseGroup(searchAlbumInfo, releaseGroup)
	releases = releasesFromReleaseGroup

	if getReleasesErr != nil {
		return release, extraReleases, releaseGrouperr
	}

	for _, releaseGroup := range otherReleaseGroups {
		releasesFromReleaseGroup, getReleasesErr = getReleasesFromReleaseGroup(searchAlbumInfo, releaseGroup)
		if getReleasesErr != nil {
			return release, extraReleases, releaseGrouperr
		}
		releases = append(releases, releasesFromReleaseGroup...)
	}

	release = releases[0]
	extraReleases = releases[1:]

	return release, extraReleases, nil
}
