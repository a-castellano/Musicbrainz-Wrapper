# Musicbrainz-Wrapper

[Project's page](https://musicmanager.gitpages.windmaker.net/Musicbrainz-Wrapper)

[Actual Repo](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper)

 [![pipeline status](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/badges/master/pipeline.svg)](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/-/commits/master) [![coverage report](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/badges/master/coverage.svg)](https://git.windmaker.net/musicmanager/Musicbrainz-Wrapper/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=music-manager-musicbrainz-wrapper&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=music-manager-musicbrainz-wrapper)

This service retrieves Artists and Records from MusicBrainz.

## Album Search steps

MusicBrainz sotres Albums inside Release Groups. In order to find requested Album, a Release Group with the same name must be searched first.

After finding the Release Group, we must find the recordings related with that Release Group in order to read Record Tracks.

## Behavior

This service retrieves Artists and Records from [MusicBrainz API](https://musicbrainz.org/doc/MusicBrainz_API).

The service receives jobs sent from [Job Manager](https://git.windmaker.net/musicmanager/Job-Manager), and process them. For each processed job this service will generate a new job containing process status and result.

### Config example

This service will look for its config in **/etc/music-manager-service/config.toml**, parent folder can be changed setting the environment variable **MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION**.

Here is a config example:

```toml
[server]
host = "localhost"
port = 5672
user = "guest"
password = "pass"

[incoming]
name = "incoming"

[outgoing]
name = "outgoing"
```

## Testing

### Unit tests

Use make to run unit tests:
```bash
make test
```

### Integration tests

Docker is required for make the following tests run, user must be sudoer too.
```bash
bash scripts/start_rabbitmq_test_server.sh
make test_integration
```
