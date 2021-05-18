// +build integration_tests unit_tests

package albums

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockOneAlbum struct {
}

func (m MockOneAlbum) SearchReleaseGroups(req *http.Request) (*http.Response, error) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-05-18T17:40:07.762Z","count":1,"offset":0,"release-groups":[{"id":"4a9421e6-9def-4616-82bf-674d5a5c7a29","type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","score":100,"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","count":1,"title":"Hexndeifl","first-release-date":"2021-02-15","primary-type":"Album","artist-credit":[{"name":"Gråinheim","artist":{"id":"c46ed989-cb1c-48fb-848d-0078784fc3f4","name":"Gråinheim","sort-name":"Gråinheim"}}],"releases":[{"id":"cb36e0ab-6634-4d85-a84e-8be89924811f","status-id":"4e304316-386d-3409-af2e-78857eec5cfe","title":"Hexndeifl","status":"Official"}]}]}
	`))}}}

	response, responseError := client.Do(req)

	return response, responseError
}

func (m MockOneAlbum) GetReleases(req *http.Request) (*http.Response, error) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"releases":[{"date":"2021-02-15","country":"XW","status-id":"4e304316-386d-3409-af2e-78857eec5cfe","text-representation":{"script":"Latn","language":"eng"},"title":"Hexndeifl","disambiguation":"","packaging-id":"119eba76-b343-3e02-a292-f0f00644bb9b","quality":"normal","packaging":"None","release-events":[{"area":{"name":"[Worldwide]","type-id":null,"sort-name":"[Worldwide]","iso-3166-1-codes":["XW"],"type":null,"id":"525d4e18-3d00-31b9-a58b-a146a916de8f","disambiguation":""},"date":"2021-02-15"}],"id":"cb36e0ab-6634-4d85-a84e-8be89924811f","barcode":null,"status":"Official"}],"primary-type-id":"f529b476-6e62-324f-b0aa-1f3e33d313fc","disambiguation":"","title":"Hexndeifl","secondary-types":[],"primary-type":"Album","id":"4a9421e6-9def-4616-82bf-674d5a5c7a29","first-release-date":"2021-02-15","secondary-type-ids":[]}
	`))}}}

	response, responseError := client.Do(req)

	return response, responseError
}

func (m MockOneAlbum) GetReleaseInfo(req *http.Request) (*http.Response, error) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"date":"2021-02-15","asin":null,"country":"XW","text-representation":{"script":"Latn","language":"eng"},"status-id":"4e304316-386d-3409-af2e-78857eec5cfe","title":"Hexndeifl","disambiguation":"","packaging-id":"119eba76-b343-3e02-a292-f0f00644bb9b","quality":"normal","media":[{"format":"Digital Media","title":"","format-id":"907a28d9-b3b2-3ef6-89a8-7b18d91d4794","track-count":6,"tracks":[{"length":336562,"position":1,"number":"1","title":"Transilvania","id":"befa427c-5f46-4cb8-84aa-ce8ddbbafbe1","recording":{"length":336562,"video":false,"first-release-date":"2021-02-15","id":"60ab5ee3-c283-467c-ace3-d07927ae166c","disambiguation":"","title":"Transilvania"}},{"number":"2","length":258088,"position":2,"title":"Vlad Draculea","id":"8bd83c92-9d96-41fa-8478-027d5c08fba3","recording":{"length":258088,"video":false,"first-release-date":"2021-02-15","id":"01b07b1d-6cff-4cbb-9c15-79afab66caf6","title":"Vlad Draculea","disambiguation":""}},{"number":"3","position":3,"length":333366,"recording":{"disambiguation":"","title":"Finstersucht","id":"baecb7a6-e05b-48b8-ae3b-300750c3f58b","first-release-date":"2021-02-15","video":false,"length":333366},"id":"ade87fb9-6f87-4635-9736-0f3b4a416550","title":"Finstersucht"},{"recording":{"length":327267,"first-release-date":"2021-02-15","video":false,"id":"651d8ac5-c710-4fbc-965c-52e693d1f75d","disambiguation":"","title":"Aetherbrand"},"id":"74ac1061-52cf-4c53-b84c-9da690c7945e","title":"Aetherbrand","number":"4","position":4,"length":327267},{"number":"5","length":393438,"position":5,"title":"Hexndeifl","id":"17660646-6191-40f7-b5c5-1b39bd4a9333","recording":{"title":"Hexndeifl","disambiguation":"","id":"478c74ce-5c6e-4a3e-983c-b3f37b998b64","first-release-date":"2021-02-15","video":false,"length":393438}},{"number":"6","length":156763,"position":6,"title":"Ausklang","id":"2246039e-315a-4ce4-bd86-6fc81338f9a5","recording":{"id":"5092d660-9fdf-485e-94ce-272fdcbf7afa","disambiguation":"","title":"Ausklang","length":156763,"video":false,"first-release-date":"2021-02-15"}}],"track-offset":0,"position":1}],"packaging":"None","release-events":[{"area":{"type":null,"disambiguation":"","id":"525d4e18-3d00-31b9-a58b-a146a916de8f","iso-3166-1-codes":["XW"],"sort-name":"[Worldwide]","type-id":null,"name":"[Worldwide]"},"date":"2021-02-15"}],"barcode":null,"id":"cb36e0ab-6634-4d85-a84e-8be89924811f","status":"Official","cover-art-archive":{"front":true,"artwork":true,"back":false,"darkened":false,"count":1}}
	`))}}}

	response, responseError := client.Do(req)

	return response, responseError
}

func TestSearchReleaseGroupWithNoResults(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-17T11:29:49.260Z","count":0,"offset":0,"release-groups":[]}
	`))}}}

	searchAlbumInfo := SearchAlbumInfo{Client: client}
	_, _, err := SearchAlbum(searchAlbumInfo, "AnyNonExistAlbum")

	if err == nil {
		t.Errorf("TestSearchReleaseGroupWithNoResult should fail.")
	}

	if err.Error() != "No release group was found." {
		t.Errorf("Error should be 'No release group was found.', not '%s'.", err)
	}

}

func TestSearchReleaseGroupWithOnlyOneResult(t *testing.T) {

	searchAlbumInfo := MockOneAlbum{}
	_, _, err := getReleaseGroup(searchAlbumInfo, "AnyNonExistAlbum", "AnyNonExistAlbum")

	if err != nil {
		t.Errorf("TestSearchReleaseGroupWithOnlyOneResult should not fail.")
	}

}
