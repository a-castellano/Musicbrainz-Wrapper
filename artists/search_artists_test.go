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

func TestSearchArtistWithOneResult(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:59:46.493Z","count":6,"offset":0,"artists":[{"id":"00eeed6b-5897-4359-8347-b8cd28375331","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":100,"name":"Manowar","sort-name":"Manowar","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"isnis":["0000000115122690"],"life-span":{"begin":"1980","ended":null},"aliases":[{"sort-name":"Man o' War","name":"Man o' War","locale":null,"type":null,"primary":null,"begin-date":null,"end-date":null}],"tags":[{"count":1,"name":"honor"},{"count":1,"name":"war"},{"count":1,"name":"norse mythology"},{"count":1,"name":"hard rock"},{"count":7,"name":"heavy metal"},{"count":1,"name":"glory"},{"count":1,"name":"symphonic metal"},{"count":1,"name":"power metal"}]},{"id":"2f3d8c8b-cd65-49e2-9a02-e89764411f88","score":39,"name":"Womanowar","sort-name":"Womanowar","life-span":{"ended":null}},{"id":"44f82b18-151f-4e98-a6f1-33b18ad6a46f","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":39,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Josef Manowarda","sort-name":"Manowarda, Josef","gender":"male","country":"DE","area":{"id":"85752fda-13c4-31a3-bee5-0e5cb1f51dad","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Germany","sort-name":"Germany","life-span":{"ended":null}},"life-span":{"begin":"1890-07-03","end":"1942-12-24","ended":true}},{"id":"4c96a93c-85ab-4aaa-83f2-bd1ad759e63b","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":38,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Eric Adams","sort-name":"Adams, Eric","gender":"male","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"disambiguation":"Louis Marullo, Manowar","life-span":{"begin":"1954-07-12","ended":null},"aliases":[{"sort-name":"Marullo, Louis","type-id":"d4dcd0c0-b341-3612-a332-c0ce797b25cf","name":"Louis Marullo","locale":null,"type":"Legal name","primary":null,"begin-date":null,"end-date":null}]},{"id":"963474e4-c3af-4d58-acbd-36a4b13f001c","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Womenowar","sort-name":"Womenowar","area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"begin-area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"disambiguation":"UK female Manowar tribute","life-span":{"begin":"2018","ended":null}},{"id":"fb1125a1-47b9-4c38-9944-4374ae785ed1","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Hanowar","sort-name":"Hanowar","disambiguation":"Bad UK Manowar tribute","life-span":{"ended":null}}]}
	`))}}}

	artistData, artistExtraData, err := SearchArtist(client, "Manowar")

	if err != nil {
		t.Errorf("TestSearchArtistWithOneResult shouldnt fail.")
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
	if artistData.Country != "US" {
		t.Errorf("Artist Country should be US, not %s.", artistData.Country)
	}

	if len(artistExtraData) != 0 {
		t.Errorf("There is only one band called Manowar, artistExtraData should be empty.")
	}
}

func TestSearchArtistWithMoreThanOneResult(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-05T11:18:27.072Z","count":48,"offset":0,"artists":[{"id":"3f03ad4f-99e1-4e52-8eb8-eb79f2a0c55e","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":100,"name":"Solstice","sort-name":"Solstice","country":"GB","area":{"id":"8a754a16-0027-3a29-b6d7-2b40ea0481ed","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United Kingdom","sort-name":"United Kingdom","life-span":{"ended":null}},"begin-area":{"id":"6d63cf38-80be-41f1-b81d-713d4763c542","type":"Subdivision","type-id":"fd3d44c5-80a1-3842-9745-2c4972d35afa","name":"Bradford","sort-name":"Bradford","life-span":{"ended":null}},"disambiguation":"epic doom metal band from United Kingdom","isnis":["0000000119412185"],"life-span":{"begin":"1990","ended":null},"tags":[{"count":2,"name":"epic doom metal"},{"count":2,"name":"doom metal"},{"count":1,"name":"metal"},{"count":1,"name":"british"},{"count":1,"name":"british doom metal"},{"count":1,"name":"uk"},{"count":2,"name":"heavy metal"},{"count":1,"name":"epic metal"}]},{"id":"06942e9b-567b-45c4-97a6-9d98265ba479","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":96,"name":"Solstice","sort-name":"Solstice","country":"GB","area":{"id":"8a754a16-0027-3a29-b6d7-2b40ea0481ed","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United Kingdom","sort-name":"United Kingdom","life-span":{"ended":null}},"begin-area":{"id":"f5923e91-2e33-423d-b276-e7df71fdd712","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Milton Keynes","sort-name":"Milton Keynes","life-span":{"ended":null}},"disambiguation":"UK neo-progressive band","isnis":["0000000119412185"],"life-span":{"begin":"1980","ended":null},"tags":[{"count":1,"name":"progressive rock"}]},{"id":"1a0fa8d8-fb02-4172-8a5c-5b0498bf7769","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":96,"name":"Solstice","sort-name":"Solstice","area":{"id":"d416a887-d13a-4ff6-8000-567ae40acca9","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Neuf-Mesnil","sort-name":"Neuf-Mesnil","life-span":{"ended":null}},"begin-area":{"id":"d416a887-d13a-4ff6-8000-567ae40acca9","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Neuf-Mesnil","sort-name":"Neuf-Mesnil","life-span":{"ended":null}},"disambiguation":"French Alternative Rock Band","life-span":{"begin":"2013-12-21","ended":null}},{"id":"ab492012-e39b-49b6-9615-840726a98eb7","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":94,"name":"Solstice Ensemble","sort-name":"Solstice Ensemble","country":"BE","area":{"id":"5b8a5ee5-0bb3-34cf-9a75-c27c44e341fc","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Belgium","sort-name":"Belgium","life-span":{"ended":null}},"disambiguation":"Belgian baroque ensemble led by Isabelle Lamfalussy","life-span":{"begin":"1997","end":"2012","ended":true},"aliases":[{"sort-name":"Ensemble Solstice","type-id":"1937e404-b981-3cb7-8151-4c86ebfc8d8e","name":"Ensemble Solstice","locale":null,"type":"Search hint","primary":null,"begin-date":null,"end-date":null},{"sort-name":"Solstice","type-id":"894afba6-2816-3c24-8072-eadb66bd04bc","name":"Solstice","locale":null,"type":"Artist name","primary":null,"begin-date":null,"end-date":null}]},{"id":"bbb8df41-678f-4bba-8814-8702c6fa78ee","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":94,"name":"Solstice","sort-name":"Solstice","disambiguation":"Celtic band from Montreal","life-span":{"ended":null}},{"id":"5d1ed9b1-1bb2-4461-a89b-761eb695558a","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":94,"name":"Solstice","sort-name":"Solstice","area":{"id":"4a9aeb42-3763-4234-8fb8-1167ac1dfdfe","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Miami","sort-name":"Miami","life-span":{"ended":null}},"disambiguation":"American metal band","life-span":{"begin":"1990","ended":null},"tags":[{"count":1,"name":"death metal"},{"count":1,"name":"thrash metal"}]},{"id":"c6b67d9d-2365-47b8-bce2-6011971dd19c","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":93,"name":"Solstice","sort-name":"Solstice","country":"AU","area":{"id":"106e0bec-b638-3b37-b731-f53d507dc00e","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Australia","sort-name":"Australia","life-span":{"ended":null}},"disambiguation":"Australian Jazz Trio","life-span":{"ended":null},"tags":[{"count":1,"name":"jazz"}]},{"id":"e0343f45-4293-4e60-8653-8ba3c1061722","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":93,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Solstice","sort-name":"Solstice","gender":"male","begin-area":{"id":"471c46a7-afc5-31c4-923c-d0444f5053a4","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Spain","sort-name":"Spain","life-span":{"ended":null}},"disambiguation":"Hardstyle by Juan Vilchez","life-span":{"ended":null}},{"id":"2e0c74c7-d1a8-49bf-b179-ea2355e8224b","score":93,"name":"Solstice","sort-name":"Solstice","country":"DE","area":{"id":"85752fda-13c4-31a3-bee5-0e5cb1f51dad","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Germany","sort-name":"Germany","life-span":{"ended":null}},"disambiguation":"Germany; released Obliquity","life-span":{"ended":null}},{"id":"1602e08c-ce1b-46f8-8f77-6551fa0a4038","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":92,"gender-id":"93452b5a-a947-30c8-934f-6a4056b151c2","name":"Solstice","sort-name":"Solstice","gender":"female","disambiguation":"Killah Priest collaborator","life-span":{"ended":null}},{"id":"f40d2ad4-187f-40e7-9a9f-6bae9f3e9286","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":92,"name":"Solstice","sort-name":"Solstice","country":"NL","area":{"id":"ef1b7cc0-cd26-36f4-8ea0-04d9623786c7","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Netherlands","sort-name":"Netherlands","life-span":{"ended":null}},"disambiguation":"Dutch Death Metal","life-span":{"ended":null}},{"id":"224b821f-b143-4d59-9f4f-4c07e954fc9c","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"appears on a dub compilation","life-span":{"ended":null}},{"id":"4bff7038-424f-4065-9c5b-c1117c4a68c9","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":92,"name":"Solstice","sort-name":"Solstice","country":"BE","area":{"id":"5b8a5ee5-0bb3-34cf-9a75-c27c44e341fc","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Belgium","sort-name":"Belgium","life-span":{"ended":null}},"begin-area":{"id":"56beda2e-dcd5-47e8-ac57-dbe72cf0fbb4","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Hasselt","sort-name":"Hasselt","life-span":{"ended":null}},"disambiguation":"Belgian rock band","life-span":{"begin":"2012","ended":null}},{"id":"a0b0ba58-4b9d-4f13-8657-dd1c8b46e524","score":92,"name":"Solstice","sort-name":"Solstice","country":"NL","area":{"id":"ef1b7cc0-cd26-36f4-8ea0-04d9623786c7","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Netherlands","sort-name":"Netherlands","life-span":{"ended":null}},"begin-area":{"id":"ef1b7cc0-cd26-36f4-8ea0-04d9623786c7","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Netherlands","sort-name":"Netherlands","life-span":{"ended":null}},"disambiguation":"Dutch Post-Rock/Pop Noir","life-span":{"ended":null}},{"id":"989f692e-628d-4bac-92bb-d945811f75a2","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":92,"name":"Solstice","sort-name":"Solstice","area":{"id":"a510b9b1-404d-4e23-8db8-0f6585909ed8","type":"Subdivision","type-id":"fd3d44c5-80a1-3842-9745-2c4972d35afa","name":"Quebec","sort-name":"Quebec","life-span":{"ended":null}},"begin-area":{"id":"a510b9b1-404d-4e23-8db8-0f6585909ed8","type":"Subdivision","type-id":"fd3d44c5-80a1-3842-9745-2c4972d35afa","name":"Quebec","sort-name":"Quebec","life-span":{"ended":null}},"end-area":{"id":"a510b9b1-404d-4e23-8db8-0f6585909ed8","type":"Subdivision","type-id":"fd3d44c5-80a1-3842-9745-2c4972d35afa","name":"Quebec","sort-name":"Quebec","life-span":{"ended":null}},"disambiguation":"Jazz fusion Québec","life-span":{"begin":"1976","end":"1981","ended":true}},{"id":"8622a9a2-b49e-4299-8997-bc94f008d7bc","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Sorrow Of Keeper","life-span":{"ended":null}},{"id":"3bc794b4-7300-4d5e-9e53-bc55bd327ff0","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Latin Fusion","life-span":{"ended":null}},{"id":"9a6f4336-0e80-44ff-a931-a5036855b762","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Solstice in 1977","life-span":{"ended":null}},{"id":"62f08d7b-302e-4c19-9cd5-8cf2e270b80f","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released J'veux du soleil","life-span":{"ended":null}},{"id":"093936cb-70dd-487e-b572-23a73efb8a42","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"A cappella","life-span":{"ended":null}},{"id":"9e623a05-1f8e-4e92-9bc6-b8a44fec1fe0","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released split with Statement","life-span":{"ended":null}},{"id":"eca9b33b-a70b-415d-b83f-0b25ddadfef4","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Punku","life-span":{"ended":null}},{"id":"37de4a13-a7ab-4c00-bdfd-3ac64dbeb040","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Smile","life-span":{"ended":null}},{"id":"4296b09d-364a-4d6d-989a-8304c342f136","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released Blue Husky","life-span":{"ended":null}},{"id":"f0250e45-29c5-4095-b419-9ac805e6c5b9","score":92,"name":"Solstice","sort-name":"Solstice","disambiguation":"released For You","life-span":{"ended":null}}]}
	`))}}}

	artistData, artistExtraData, err := SearchArtist(client, "Solstice")

	if err != nil {
		t.Errorf("TestClientNoArtists shouldn't fail.")
	}

	if artistData.Name != "Solstice" {
		t.Errorf("Artist Name should be Solstice, not %s.", artistData.Name)
	}
	if artistData.URL != "https://musicbrainz.org/artist/3f03ad4f-99e1-4e52-8eb8-eb79f2a0c55e" {
		t.Errorf("Artist URL should be https://musicbrainz.org/artist/3f03ad4f-99e1-4e52-8eb8-eb79f2a0c55e, not %s.", artistData.URL)
	}
	if artistData.ID != "3f03ad4f-99e1-4e52-8eb8-eb79f2a0c55e" {
		t.Errorf("Artist ID should be 3f03ad4f-99e1-4e52-8eb8-eb79f2a0c55e, not %s.", artistData.ID)
	}
	if artistData.Country != "GB" {
		t.Errorf("Artist Country should be GB, not %s.", artistData.Country)
	}

	if len(artistExtraData) == 0 {
		t.Errorf("There are more than one band called Solstice, artistExtraData shouldn't' be empty.")
	}

}
