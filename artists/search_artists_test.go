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

func TestSearchArtistWithMoreThanOneResult(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:59:46.493Z","count":6,"offset":0,"artists":[{"id":"00eeed6b-5897-4359-8347-b8cd28375331","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":100,"name":"Manowar","sort-name":"Manowar","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"isnis":["0000000115122690"],"life-span":{"begin":"1980","ended":null},"aliases":[{"sort-name":"Man o' War","name":"Man o' War","locale":null,"type":null,"primary":null,"begin-date":null,"end-date":null}],"tags":[{"count":1,"name":"honor"},{"count":1,"name":"war"},{"count":1,"name":"norse mythology"},{"count":1,"name":"hard rock"},{"count":7,"name":"heavy metal"},{"count":1,"name":"glory"},{"count":1,"name":"symphonic metal"},{"count":1,"name":"power metal"}]},{"id":"2f3d8c8b-cd65-49e2-9a02-e89764411f88","score":39,"name":"Womanowar","sort-name":"Womanowar","life-span":{"ended":null}},{"id":"44f82b18-151f-4e98-a6f1-33b18ad6a46f","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":39,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Josef Manowarda","sort-name":"Manowarda, Josef","gender":"male","country":"DE","area":{"id":"85752fda-13c4-31a3-bee5-0e5cb1f51dad","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Germany","sort-name":"Germany","life-span":{"ended":null}},"life-span":{"begin":"1890-07-03","end":"1942-12-24","ended":true}},{"id":"4c96a93c-85ab-4aaa-83f2-bd1ad759e63b","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":38,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Eric Adams","sort-name":"Adams, Eric","gender":"male","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"disambiguation":"Louis Marullo, Manowar","life-span":{"begin":"1954-07-12","ended":null},"aliases":[{"sort-name":"Marullo, Louis","type-id":"d4dcd0c0-b341-3612-a332-c0ce797b25cf","name":"Louis Marullo","locale":null,"type":"Legal name","primary":null,"begin-date":null,"end-date":null}]},{"id":"963474e4-c3af-4d58-acbd-36a4b13f001c","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Womenowar","sort-name":"Womenowar","area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"begin-area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"disambiguation":"UK female Manowar tribute","life-span":{"begin":"2018","ended":null}},{"id":"fb1125a1-47b9-4c38-9944-4374ae785ed1","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Hanowar","sort-name":"Hanowar","disambiguation":"Bad UK Manowar tribute","life-span":{"ended":null}}]}
	`))}}}

	artistData, _, err := SearchArtist(client, "Manowar")

	if err != nil {
		t.Errorf("TestClientNoArtists yyshouldnt fail.")
	}

	if artistData.Name != "Manowar" {
		t.Errorf("Artist Name should be Manowar, not %s.", artistData.Name)
	}

	if artistData.URL != "https://musicbrainz.org/artist/00eeed6b-5897-4359-8347-b8cd28375331" {
		t.Errorf("Artist URL should be https://musicbrainz.org/artist/00eeed6b-5897-4359-8347-b8cd28375331, not %s.", artistData.URL)
	}
	if artistData.ID != "00eeed6b-5897-4359-8347-b8cd28375331" {
		t.Errorf("Artist ID should be 00eeed6b-5897-4359-8347-b8cd28375331, not %s.", artistData.ID)
	}
}
