package artists

import (
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"net/http"
)

func GetArtistRecords(client http.Client, artistData SearchArtistData) ([]commontypes.Record, error) {
	var records []commontypes.Record
	return records, nil
}
