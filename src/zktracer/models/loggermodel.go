package models

type LoggerModel struct {
	Id     string
	Addr   string
	BlkDev string
	Status bool
	Cpu    float64
	Net    float64
	Disk   float64
}
