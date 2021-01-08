// +build integration_tests unit_tests

package artists

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetArtistRecordBroken(t *testing.T) {
	artistData := SearchArtistData{Name: "Tombs", URL: "https://musicbrainz.org/artist/8b2cccef-ad8b-4029-911e-14157e4cd5ae", ID: "8b2cccef-ad8b-4029-911e-14157e4cd5ae", Genre: "", Country: "US"}
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
non valid json
`))}}}

	records, err := GetArtistRecords(client, artistData)

	if err == nil {
		t.Errorf("TestGetArtistRecordBroken should fail.")
	}

	if len(records) != 0 {
		t.Errorf("Number of retrieved records length in TestGetArtistRecordBroken should be 0.")
	}
}

func TestGetArtistRecords(t *testing.T) {
	artistData := SearchArtistData{Name: "Tombs", URL: "https://musicbrainz.org/artist/8b2cccef-ad8b-4029-911e-14157e4cd5ae", ID: "8b2cccef-ad8b-4029-911e-14157e4cd5ae", Genre: "", Country: "US"}
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"disambiguation":"","isnis":[],"country":"US","life-span":{"ended":false,"end":null,"begin":"2007"},"gender-id":null,"id":"8b2cccef-ad8b-4029-911e-14157e4cd5ae","sort-name":"Tombs","type":"Group","gender":null,"end_area":null,"begin_area":{"type-id":null,"sort-name":"Brooklyn","type":null,"name":"Brooklyn","id":"a71b0d32-7752-49e9-8594-2247ad6ac12c","disambiguation":""},"area":{"id":"489ce91b-6658-3307-9877-795b68554c98","disambiguation":"","iso-3166-1-codes":["US"],"type-id":null,"type":null,"name":"United States","sort-name":"United States"},"ipis":[],"begin-area":{"type-id":null,"sort-name":"Brooklyn","type":null,"name":"Brooklyn","id":"a71b0d32-7752-49e9-8594-2247ad6ac12c","disambiguation":""},"end-area":null,"name":"Tombs","release-groups":[{"secondary-type-ids":[],"first-release-date":"2009-04-17","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","title":"Winterhours","secondary-types":[],"disambiguation":"","primary-type":"Album","id":"6864f2f5-8609-3abc-857c-c143da8e1723"},{"disambiguation":"","primary-type":"Album","id":"8cc5ccd4-7d9f-4856-b39d-7228f8ac77a3","title":"Path of Totality","secondary-types":[],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","secondary-type-ids":[],"first-release-date":"2011-06-07"},{"first-release-date":"2012-05-15","secondary-type-ids":[],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","secondary-types":[],"title":"Label Showcase - Relapse Records","id":"cf641763-a690-4e7e-b875-eae8298163fc","primary-type":"Album","disambiguation":""},{"title":"Savage Gold","secondary-types":[],"disambiguation":"","primary-type":"Album","id":"d090e390-cbf8-4d79-b8df-b3e5d3a0d8c0","secondary-type-ids":[],"first-release-date":"2014-06-10","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc"},{"first-release-date":"2017-06-16","secondary-type-ids":[],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","secondary-types":[],"title":"The Grand Annihilation","id":"1aecb1cf-0758-49f0-b922-e54ec264833c","primary-type":"Album","disambiguation":""},{"secondary-type-ids":[],"first-release-date":"2020-11-20","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","title":"Under Sullen Skies","secondary-types":[],"disambiguation":"","primary-type":"Album","id":"00e8173f-32ed-48bb-8d64-d58d7f61cf7a"},{"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","secondary-type-ids":["dd2a21e1-0c00-3729-a7a0-de60b84eb5d1"],"first-release-date":"2010","primary-type":"Album","disambiguation":"","id":"6084e220-e460-41d1-853f-858d4002f75b","secondary-types":["Compilation"],"title":"Fear Is the Weapon"},{"title":"Tombs / Planks","secondary-types":[],"id":"8dc65a7d-2ab6-4d92-942d-6ca5469c1ab5","disambiguation":"","primary-type":"Single","first-release-date":"2008","secondary-type-ids":[],"primary-type-id":"d6038452-8ee0-3f68-affc-2de9a1ede0b9"},{"primary-type-id":"d6038452-8ee0-3f68-affc-2de9a1ede0b9","first-release-date":"2012-09","secondary-type-ids":[],"id":"f6e6f89e-f004-497d-a505-70c504a54afc","primary-type":"Single","disambiguation":"","secondary-types":[],"title":"Ashes"},{"disambiguation":"","primary-type":"EP","id":"4bb63d88-054e-3328-8baf-35f5a50365fa","title":"Tombs","secondary-types":[],"primary-type-id":"6d0c5bf6-7a33-3420-a519-44fc63eedebf","secondary-type-ids":[],"first-release-date":"2007"},{"secondary-types":[],"title":"All Empires Fall","primary-type":"EP","disambiguation":"","id":"59f51b72-6a76-41d0-847c-e9a90bac4f44","secondary-type-ids":[],"first-release-date":"2016-04-01","primary-type-id":"6d0c5bf6-7a33-3420-a519-44fc63eedebf"},{"secondary-type-ids":[],"first-release-date":"2020-02-28","primary-type-id":"6d0c5bf6-7a33-3420-a519-44fc63eedebf","title":"Monarchy of Shadows","secondary-types":[],"disambiguation":"","primary-type":"EP","id":"d4684fc1-26f7-437d-b0f3-63ab300368f1"}],"type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b"}
`))}}}

	records, err := GetArtistRecords(client, artistData)

	if err != nil {
		t.Errorf("TestGetArtistRecords shouldn't fail.")
	}

	if len(records) != 12 {
		t.Errorf("Number of retrieved records length should be 6.")
	}
}
