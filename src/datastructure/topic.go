package datastructure

type TopicInfo struct {
	AppId      string    `json:"app_id"`
	Name       string    `json:"topic_name"`
	ReplicaNum int       `json:"replica_num"`
	Retention  int       `json:"retention"`
	Segments   []Segment `json:"segments"`
}

type Segment struct {
	LastConfirmEntry uint64   `json:"last_confirm_entry"`
	Loggers          []string `json:"loggers"`
}
