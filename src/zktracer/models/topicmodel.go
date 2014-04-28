package models

type TopicModel struct {
	Name       string
	AppId      string
	BrokerId   string
	ReplicaNum int
	Retention  int
	Segments   string
	Status     bool
}
