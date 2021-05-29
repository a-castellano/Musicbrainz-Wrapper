# Musicbrainz-Wrapper

[Project's page](https://musicmanager.gitpages.windmaker.net/Musicbrainz-Wrapper)

[Actual Repo](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper)

 [![pipeline status](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/badges/master/pipeline.svg)](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/-/commits/master) [![coverage report](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/badges/master/coverage.svg)](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=music-manager-musicbrainz-wrapper&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=music-manager-musicbrainz-wrapper)

This service retrieves Artists and Records from MusicBrainz.

## Album Search steps

MusicBrainz sotres Albums inside Release Groups. In order to find requested Album, a Release Group with the same name must be searched first.

After finding the Release Group, we must find the recordings related with that Release Group in order to read Record Tracks.

