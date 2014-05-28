package main

import (
	//"datastructure"
	"encoding/json"
	"fmt"
	//"math/rand"
	//"strconv"
	"sync"
	//"sysutil"
	"time"
	"zk"
)

// zones：zoneid对应的所有主机ip
// rzs：相对于zoneid来说的远程zone集合
// remote_zones：远程zone的集合
var (
	zones = map[string][]string{
		//"zone0": []string{"172.19.32.16"},
		"zone1": []string{"172.19.32.16"},
		"zone2": []string{"172.19.32.153"},
		"zone3": []string{"172.19.32.46"},
	}
	rzs          = []string{}
	remote_zones = make(map[string][]string)
)

// 生成zoneid到remote_zones的映射
func GetRemoteZonesMap() {
	// 生成rzs集合
	for k, _ := range zones {
		rzs = append(rzs, k)
	}
	// 生成remote_zones映射
	for index, value := range rzs {
		remote_zones[value] = DeleteElementInSlice(rzs, index)
	}
}

// 从slice中删除指定元素
func DeleteElementInSlice(s []string, index int) []string {
	ret := make([]string, len(s))
	copy(ret, s)
	ret = append(ret[:index], ret[index+1:]...)
	return ret
}

func main() {
	fmt.Printf("Data emulator started.\n")
	GetRemoteZonesMap()
	// 加锁等所有协程结束后再退出程序
	var wg sync.WaitGroup
	for _, e := range rzs {
		wg.Add(1)
		go func(zoneid string) {
			defer wg.Done()
			Work(zoneid)
		}(e)
	}
	wg.Wait()
}

func Work(zoneid string) {
	conn := zk.Connect(zones[zoneid], 1*time.Second)

	// 先清空原来的数据
	err := conn.DeleteRecur("/ymb")
	if err != nil {
		panic(err)
	}
	// 加入【/ymb】节点
	if flag, err := conn.Exists("/ymb"); err == nil && !flag {
		conn.Create("/ymb", "", zk.WorldACL(zk.PermAll), 0)
	}
	// 加入【/ymb/zoneid】节点
	if flag, err := conn.Exists("/ymb/zoneid"); err == nil && !flag {
		conn.Create("/ymb/zoneid", zoneid, zk.WorldACL(zk.PermAll), 0)
	}
	// 加入【/ymb/remote_zones】节点
	if flag, err := conn.Exists("/ymb/remote_zones"); err == nil && !flag {
		remote_zones_json, err := json.Marshal(remote_zones[zoneid])
		if err != nil {
			panic(err)
		}
		conn.Create("/ymb/remote_zones", string(remote_zones_json), zk.WorldACL(zk.PermAll), 0)
	}
	// 加入【/ymb/brokers】节点
	if flag, err := conn.Exists("/ymb/brokers"); err == nil && !flag {
		conn.Create("/ymb/brokers", "", zk.WorldACL(zk.PermAll), 0)
	}
	// 加入【/ymb/loggers】节点
	if flag, err := conn.Exists("/ymb/loggers"); err == nil && !flag {
		conn.Create("/ymb/loggers", "", zk.WorldACL(zk.PermAll), 0)
	}
	// 加入【/ymb/loggers】节点
	if flag, err := conn.Exists("/ymb/loggers"); err == nil && !flag {
		//conn.Create("/ymb/loggers", "", zk.WorldACL(zk.PermAll), 0)
	}

	// 每10秒刷新一下数据
	timer := time.Tick(2 * time.Second)
	for _ = range timer {
		// 生成100个broker
		// 加入【/ymb/brokers】节点
		// 加入【/ymb/loggers】节点
		// 加入【/ymb/topics】节点
		// 加入【/ymb/topics/appid/topic.name/broker】节点
		// 加入【/ymb/appid】节点

		flag, err := conn.Exists("/")
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%s]: %v\n", zoneid, flag)
	}

	conn.Close()
}
