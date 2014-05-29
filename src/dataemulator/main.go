package main

import (
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/xianmau/gozk/zk"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sysutil"
	"time"
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

var (
	// 假设是电信的IP集合
	ipset_1 = []string{
		"183.63.54.99",
		"183.63.36.57",
		"183.63.42.18",
		"183.63.101.29",
		"113.95.87.61",
		"113.95.91.10",
		"113.95.43.33",
		"219.143.58.107",
		"219.143.36.24",
		"219.143.39.83",
	}
	// 假设是联通的IP集合
	ipset_2 = []string{
		"120.87.14.87",
		"120.87.66.13",
		"120.87.97.101",
		"120.87.120.240",
		"163.125.8.45",
		"163.125.63.14",
		"163.125.12.88",
		"218.104.56.40",
		"218.104.32.41",
		"218.104.96.138",
	}
	// 假设是磁盘块的名称集合
	bledevset = []string{
		"sda",
		"sdb",
		"sdc",
		"sdd",
		"sde",
	}
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
			Work(zoneid)
			wg.Done()
		}(e)
	}
	wg.Wait()
}

func Work(zoneid string) {
	// 创建随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	conn := zk.Connect(zones[zoneid], 1*time.Second)
	// 先清空原来的数据
	err := conn.DeleteRecur("/ymb")
	if err != nil {
		panic(err)
	}

	//加入【/ymb】节点
	if flag, err := conn.Exists("/ymb"); err == nil && !flag {
		conn.Create("/ymb", "", zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb")
	}
	// 加入【/ymb/zoneid】节点
	if flag, err := conn.Exists("/ymb/zoneid"); err == nil && !flag {
		conn.Create("/ymb/zoneid", zoneid, zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/zoneid")
	}
	// 加入【/ymb/remote_zones】节点
	if flag, err := conn.Exists("/ymb/remote_zones"); err == nil && !flag {
		remote_zones_json, err := json.Marshal(remote_zones[zoneid])
		if err != nil {
			panic(err)
		}
		conn.Create("/ymb/remote_zones", string(remote_zones_json), zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/remote_zones")
	}
	// 加入【/ymb/brokers】节点
	if flag, err := conn.Exists("/ymb/brokers"); err == nil && !flag {
		conn.Create("/ymb/brokers", "", zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/brokers")
	}
	// 加入【/ymb/loggers】节点
	if flag, err := conn.Exists("/ymb/loggers"); err == nil && !flag {
		conn.Create("/ymb/loggers", "", zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/loggers")
	}
	// 加入【/ymb/appid】节点
	if flag, err := conn.Exists("/ymb/appid"); err == nil && !flag {
		conn.Create("/ymb/appid", "", zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/appid")
	}
	// 加入【/ymb/topics】节点
	if flag, err := conn.Exists("/ymb/topics"); err == nil && !flag {
		conn.Create("/ymb/topics", "", zk.WorldACL(zk.PermAll), 0)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/topics")
	}

	// 生成app数据
	app_set := []string{}
	for i := 0; i < 10; i++ {
		app := datastructure.AppInfo{
			Id:  "app_" + strconv.Itoa(i),
			Key: "KEYIS" + strconv.Itoa(1000+r.Intn(9000)),
		}
		app_json, err := json.Marshal(app)
		if err != nil {
			panic(err)
		}
		if flag, err := conn.Exists("/ymb/appid/" + app.Id); err == nil && !flag {
			conn.Create("/ymb/appid/"+app.Id, string(app_json), zk.WorldACL(zk.PermAll), 0)
		}
		app_set = append(app_set, app.Id)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/appid/"+app.Id)
	}
	// 生成broker数据
	broker_set := []string{}
	for i := 0; i < 100; i++ {
		broker := datastructure.BrokerInfo{
			Id:    "broker_" + strconv.Itoa(i),
			Addrs: []string{ipset_1[r.Intn(len(ipset_1))], ipset_2[r.Intn(len(ipset_2))]},
			Stat:  sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()},
		}
		broker_json, err := json.Marshal(broker)
		if err != nil {
			panic(err)
		}
		if flag, err := conn.Exists("/ymb/brokers/" + broker.Id); err == nil && !flag {
			conn.Create("/ymb/brokers/"+broker.Id, string(broker_json), zk.WorldACL(zk.PermAll), 0)
		}
		broker_set = append(broker_set, broker.Id)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/brokers/"+broker.Id)
	}
	// 生成logger数据
	logger_set := []string{}
	for i := 0; i < 1000; i++ {
		logger := datastructure.LoggerInfo{
			Id:     "logger_" + strconv.Itoa(i),
			Addr:   ipset_1[r.Intn(len(ipset_1))],
			BlkDev: bledevset[r.Intn(len(bledevset))],
			Stat:   sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()},
		}
		logger_json, err := json.Marshal(logger)
		if err != nil {
			panic(err)
		}
		if flag, err := conn.Exists("/ymb/loggers/" + logger.Id); err == nil && !flag {
			conn.Create("/ymb/loggers/"+logger.Id, string(logger_json), zk.WorldACL(zk.PermAll), 0)
		}
		logger_set = append(logger_set, logger.Id)
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/loggers/"+logger.Id)
	}
	// 生成topic数据
	for i := 0; i < 10000; i++ {
		topic := datastructure.TopicInfo{
			AppId:      app_set[r.Intn(len(app_set))],
			Name:       "NAMEIS" + strconv.Itoa(100000+r.Intn(900000)),
			ReplicaNum: 3,
			Retention:  60,
			Segments: []datastructure.Segment{
				datastructure.Segment{uint64(r.Int63()), []string{"logger1", "logger2"}},
				datastructure.Segment{uint64(r.Int63()), []string{"logger1", "logger2"}},
			},
		}
		brokerid := broker_set[r.Intn(len(broker_set))]
		topic_json, err := json.Marshal(topic)
		if err != nil {
			panic(err)
		}
		if flag, err := conn.Exists("/ymb/topics/" + topic.AppId); err == nil && !flag {
			conn.Create("/ymb/topics/"+topic.AppId, "", zk.WorldACL(zk.PermAll), 0)
		}
		if flag, err := conn.Exists("/ymb/topics/" + topic.AppId + "/" + topic.Name); err == nil && !flag {
			conn.Create("/ymb/topics/"+topic.AppId+"/"+topic.Name, string(topic_json), zk.WorldACL(zk.PermAll), 0)
		}
		if flag, err := conn.Exists("/ymb/topics/" + topic.AppId + "/" + topic.Name + "/" + brokerid); err == nil && !flag {
			conn.Create("/ymb/topics/"+topic.AppId+"/"+topic.Name+"/"+brokerid, "", zk.WorldACL(zk.PermAll), 0)
		}
		log.Printf("%s: node [%s] created\n", zoneid, "/ymb/topics/"+topic.AppId+"/"+topic.Name+"/"+brokerid)
	}

	// 每10秒刷新一下数据
	timer := time.Tick(2 * time.Second)
	for _ = range timer {
		// 生成broker数据
		for _, bro := range broker_set {
			// 生成path
			path := "/ymb/brokers/" + bro
			// 获取data
			broker_json, err := conn.Get(path)
			if err != nil {
				panic(err)
			}
			var broker datastructure.BrokerInfo
			err = json.Unmarshal([]byte(broker_json), &broker)
			broker.Stat = sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()}
			if err != nil {
				panic(err)
			}
			broker_json_new, err := json.Marshal(broker)
			if err != nil {
				panic(err)
			}
			conn.Set("/ymb/brokers/"+broker.Id, string(broker_json_new))
		}
		log.Printf("%s: brokers updated\n", zoneid)
		// 生成logger数据
		for _, lo := range logger_set {
			// 生成path
			path := "/ymb/loggers/" + lo
			// 获取data
			logger_json, err := conn.Get(path)
			if err != nil {
				panic(err)
			}
			var logger datastructure.LoggerInfo
			err = json.Unmarshal([]byte(logger_json), &logger)
			logger.Stat = sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()}
			if err != nil {
				panic(err)
			}
			logger_json_new, err := json.Marshal(logger)
			if err != nil {
				panic(err)
			}
			conn.Set("/ymb/loggers/"+logger.Id, string(logger_json_new))
		}
		log.Printf("%s: loggers updated\n", zoneid)
	}
	conn.Close()
}
