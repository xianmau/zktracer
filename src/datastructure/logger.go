package datastructure

import (
	"sysutil"
)

type LoggerInfo struct {
	Id     string           `json:"id"`
	Addr   string           `json:"address"`
	BlkDev string           `json:"blockdev"`
	Stat   sysutil.SysUtils `json:"sysutil"`
}
