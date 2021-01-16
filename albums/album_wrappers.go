package albums

type ReleaseGroupWrapperInterface interface {
	GetReleaseRecords(url string) map[string]interface{}
}

type ReleaseWrapperInterface interface {
	GetReleaseRecords(url string) map[string]interface{}
}
