package models

type BrokerModel struct {
	Id     string
	Addrs  string
	Status bool
	Cpu    float64
	Net    float64
	Disk   float64
}
