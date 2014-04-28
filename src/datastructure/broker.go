package datastructure

import (
	"sysutil"
)

type BrokerInfo struct {
	Id    string           `json:"id"`
	Addrs []string         `json:"address"`
	Stat  sysutil.SysUtils `json:"sysutil"`
}
