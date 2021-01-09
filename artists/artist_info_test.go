// +build integration_tests unit_tests

package artists

import (
	"bytes"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
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
		t.Errorf("Number of retrieved records length should be 12.")
	}

	first_record := records[0]

	if first_record.Name != "Winterhours" {
		t.Errorf("First Tombs record should be 'Winterhours'.")
	}

	if first_record.ID != "6864f2f5-8609-3abc-857c-c143da8e1723" {
		t.Errorf("First Tombs record ID should be '6864f2f5-8609-3abc-857c-c143da8e1723'.")
	}

	if first_record.URL != "https://musicbrainz.org/release-group/6864f2f5-8609-3abc-857c-c143da8e1723" {
		t.Errorf("First Tombs record URL should be 'https://musicbrainz.org/release-group/6864f2f5-8609-3abc-857c-c143da8e1723'.")
	}

	if first_record.Year != 2009 {
		t.Errorf("First Tombs record year should be 2009.")
	}

	if first_record.Type != commontypes.FullLength {
		t.Errorf("First Tombs should be FullLength.")
	}

}

func TestGetWithEpsAndOtherTypes(t *testing.T) {
	artistData := SearchArtistData{Name: "Manowar", URL: "https://musicbrainz.org/artist/00eeed6b-5897-4359-8347-b8cd28375331", ID: "00eeed6b-5897-4359-8347-b8cd28375331", Genre: "", Country: "US"}
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","gender":null,"ipis":[],"country":null,"end-area":null,"end_area":null,"sort-name":"Perkele","gender-id":null,"area":{"sort-name":"Gothenburg","type":null,"name":"Gothenburg","id":"8f6c316e-9924-48ea-967b-16757dd82399","disambiguation":"","type-id":null},"id":"9d62c5a5-a94f-4e22-b40f-e3f394191b3f","begin-area":null,"life-span":{"ended":false,"end":null,"begin":"1993"},"isnis":[],"release-groups":[{"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","disambiguation":"","id":"7be87ab7-6185-309a-9de4-5fab22b85498","primary-type":"Album","first-release-date":"1998","title":"Från flykt till kamp","secondary-type-ids":[],"secondary-types":[]},{"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","disambiguation":"","id":"6bc3b679-d8ac-382c-aadb-9c56e303d193","primary-type":"Album","first-release-date":"2001","title":"Voice of Anger","secondary-type-ids":[],"secondary-types":[]},{"disambiguation":"","id":"5ad30cca-bfc5-39e6-bf42-1271e7d004ba","primary-type":"Album","first-release-date":"2002","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","secondary-types":[],"title":"No Shame","secondary-type-ids":[]},{"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","id":"4c4bb914-fa15-3189-8e28-3845bec65763","disambiguation":"","primary-type":"Album","first-release-date":"2003","title":"Stories From the Past","secondary-type-ids":[],"secondary-types":[]},{"secondary-type-ids":[],"title":"Confront","secondary-types":[],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","first-release-date":"2005","primary-type":"Album","id":"495064c7-a65f-36f6-952d-c0990222d459","disambiguation":""},{"secondary-types":[],"secondary-type-ids":[],"title":"Längtan","first-release-date":"2008-12-01","id":"fc003b8f-57b9-4881-8570-3ab1fed3b4f9","primary-type":"Album","disambiguation":"","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc"},{"secondary-type-ids":[],"title":"Perkele Forever","secondary-types":[],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","first-release-date":"2010","id":"543cf854-b36b-4c76-91bb-63a54a4cb404","primary-type":"Album","disambiguation":""},{"secondary-types":[],"title":"A Way Out","secondary-type-ids":[],"id":"cd779644-abff-43cc-b5b8-03e45c7371cd","primary-type":"Album","disambiguation":"","first-release-date":"2013-10-04","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc"},{"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","disambiguation":"","id":"4d116a2c-d79e-4143-b1a8-47a9030b132f","primary-type":"Album","first-release-date":"2019","title":"Leaders Of Tomorrow","secondary-type-ids":[],"secondary-types":[]},{"secondary-type-ids":["dd2a21e1-0c00-3729-a7a0-de60b84eb5d1"],"title":"Det Var Då","secondary-types":["Compilation"],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","first-release-date":"2013","id":"7d7f1aaf-9fdc-462f-b71f-30ba00d44cf8","primary-type":"Album","disambiguation":""},{"secondary-types":["Live"],"secondary-type-ids":["6fd474e2-6b58-3102-9d17-d6f7eb7da0a0"],"title":"Songs for You","first-release-date":"2009","disambiguation":"","id":"feee336f-347f-40e1-b8e1-6800d9a1a543","primary-type":"Album","primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc"},{"title":"Days of Punk","secondary-type-ids":[],"secondary-types":[],"primary-type-id":"6d0c5bf6-7a33-3420-a519-44fc63eedebf","id":"b25f2142-a09f-3a19-b714-b304a320e686","primary-type":"EP","disambiguation":"","first-release-date":"2003"},{"primary-type":"EP","id":"1c177b32-32de-3cd0-9883-f4295d6035dc","disambiguation":"","first-release-date":"2003","primary-type-id":"6d0c5bf6-7a33-3420-a519-44fc63eedebf","secondary-types":[],"title":"Göteborg","secondary-type-ids":[]},{"primary-type-id":null,"first-release-date":"2016-12-02","id":"5d812991-3dc0-493b-8c46-c30a38a65ad8","disambiguation":"","primary-type":null,"secondary-type-ids":["dd2a21e1-0c00-3729-a7a0-de60b84eb5d1"],"title":"Best From The Past","secondary-types":["Compilation"]},{"secondary-types":["Demo"],"title":"Perkele","secondary-type-ids":["81598169-0d6c-3bce-b4be-866fa658eda3"],"id":"9a01f10f-3504-4191-adcc-12727948e6ad","primary-type":null,"disambiguation":"","first-release-date":"1994","primary-type-id":null}],"begin_area":null,"type":"Group","disambiguation":"","name":"Perkele"}
`))}}}

	records, err := GetArtistRecords(client, artistData)

	if err != nil {
		t.Errorf("TestGetArtistRecords shouldn't fail.")
	}

	if len(records) != 15 {
		t.Errorf("Number of retrieved records length should be 15, not %d.", len(records))
	}

	compilation_record := records[9]

	if compilation_record.Name != "Det Var Då" {
		t.Errorf("Perkele compilation record should be 'Det Var Då', not '%s'", compilation_record.Name)
	}

	if compilation_record.ID != "7d7f1aaf-9fdc-462f-b71f-30ba00d44cf8" {
		t.Errorf("Perkele compilation record ID should be '7d7f1aaf-9fdc-462f-b71f-30ba00d44cf8', not '%s'", compilation_record.ID)
	}

	if compilation_record.Year != 2013 {
		t.Errorf("Perkele compilation record year should be 2013, not %d.", compilation_record.Year)
	}

	if compilation_record.Type != commontypes.Compilation {
		t.Errorf("Perkele compilation record should be Compilation.")
	}

	ep_record := records[12]

	if ep_record.Name != "Göteborg" {
		t.Errorf("Perkele EP record should be 'Göteborg', not '%s'", ep_record.Name)
	}

	if ep_record.ID != "1c177b32-32de-3cd0-9883-f4295d6035dc" {
		t.Errorf("Perkele EP record ID should be '1c177b32-32de-3cd0-9883-f4295d6035dc', not '%s'", ep_record.ID)
	}

	if ep_record.Year != 2003 {
		t.Errorf("Perkele EP record year should be 2003, not %d.", ep_record.Year)
	}

	if ep_record.Type != commontypes.EP {
		t.Errorf("Perkele EP record should be EP.")
	}

	live_record := records[10]

	if live_record.Name != "Songs for You" {
		t.Errorf("Perkele Live record should be 'Songs for You', not '%s'", live_record.Name)
	}

	if live_record.ID != "feee336f-347f-40e1-b8e1-6800d9a1a543" {
		t.Errorf("Perkele Live record ID should be 'feee336f-347f-40e1-b8e1-6800d9a1a543', not '%s'", live_record.ID)
	}

	if live_record.Year != 2009 {
		t.Errorf("Perkele Live record year should be 2009, not %d.", live_record.Year)
	}

	if live_record.Type != commontypes.Live {
		t.Errorf("Perkele Live record should be Live.")
	}

	demo_record := records[14]

	if demo_record.Name != "Perkele" {
		t.Errorf("Perkele Demo record should be 'Perkel', not '%s'", demo_record.Name)
	}

	if demo_record.ID != "feee336f-347f-40e1-b8e1-6800d9a1a543" {
		t.Errorf("Perkele Demo record ID should be 'feee336f-347f-40e1-b8e1-6800d9a1a543', not '%s'", demo_record.ID)
	}

	if demo_record.Year != 1994 {
		t.Errorf("Perkele Demo record year should be 1994, not %d.", demo_record.Year)
	}

	if demo_record.Type != commontypes.Demo {
		t.Errorf("Perkele demo record should be Demo.")
	}

}
