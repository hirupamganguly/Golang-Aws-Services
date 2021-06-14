package contentuploader

type ContentUploader struct {
}
type ContentUploaderRepository interface {
	Add(ei ExtractInfo) (err error)
	Get() (eis []ExtractInfo, err error)
}
type Gender int8

const (
	NotSpecified Gender = iota
	Male
	Female
)

type ExtractInfo struct {
	Name    string `json:"name" bson:"name"`
	IsTrue  bool   `json:"is_true" bson:"is_true"`
	Payload string `json:"payload" bson:"payload"`
	DataId  string `json:"data_id" bson:"data_id"`
	Version string `json:"version" bson:"version"`
	Type    string `json:"type" bson:"type"`
	Gender  Gender `json:"gender" bson:"gender"`
	Label   string `json:"label" bson:"label"`
}
