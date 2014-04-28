package main

import (
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
	"strconv"
	"sysutil"
	"time"
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

func main() {
	// 建一些数据节点
	ymb := new(datastructure.YMB)
	ymb.Apps = []datastructure.AppInfo{}
	ymb.Brokers = []datastructure.BrokerInfo{}
	ymb.Loggers = []datastructure.LoggerInfo{}
	ymb.Topics = []datastructure.TopicInfo{}
	ymb.ZoneId = string("172.19.32.35:2181")
	ymb.RemoteZones = []string{
		"172.19.32.135:2181",
		"172.19.32.124:2181",
	}

	// 创建随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成broker数据
	for i := 0; i < 100; i++ {
		var broker datastructure.BrokerInfo
		if i < 10 {
			broker.Id = "broker_00" + strconv.Itoa(i)
		} else {
			broker.Id = "broker_0" + strconv.Itoa(i)
		}
		broker.Addrs = []string{ipset_1[r.Intn(len(ipset_1))], ipset_2[r.Intn(len(ipset_2))]}
		broker.Stat = sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()}
		ymb.Brokers = append(ymb.Brokers, broker)
	}

	// 生成logger数据
	for i := 0; i < 1000; i++ {
		var logger datastructure.LoggerInfo
		if i < 10 {
			logger.Id = "logger_000" + strconv.Itoa(i)
		} else if i < 100 {
			logger.Id = "logger_00" + strconv.Itoa(i)
		} else {
			logger.Id = "logger_0" + strconv.Itoa(i)
		}
		logger.Addr = ipset_1[r.Intn(len(ipset_1))]
		logger.BlkDev = bledevset[r.Intn(len(bledevset))]
		logger.Stat = sysutil.SysUtils{r.Float64(), r.Float64(), r.Float64()}
		ymb.Loggers = append(ymb.Loggers, logger)
	}

	// 创建zk实例
	conn, _, err := zk.Connect(ymb.RemoteZones, time.Second)
	checkErr(err)
	checkErr(deleteRecur(conn, "/ymb/brokers"))
	checkErr(deleteRecur(conn, "/ymb/loggers"))
	fmt.Printf("All existing znodes deleted.\nNow create some znodes...\n")
	_, err = conn.Create("/ymb/brokers", []byte{0}, 0, zk.WorldACL(zk.PermAll))
	checkErr(err)
	_, err = conn.Create("/ymb/loggers", []byte{0}, 0, zk.WorldACL(zk.PermAll))
	checkErr(err)

	// 导入broker数据
	for _, v := range ymb.Brokers {
		buf, err := json.Marshal(v)
		checkErr(err)
		ret_p, err := conn.Create("/ymb/brokers/"+v.Id, buf, 0, zk.WorldACL(zk.PermAll))
		checkErr(err)
		fmt.Printf("ZNode[%s] created.\n", ret_p)
	}
	// 导入logger数据
	for _, v := range ymb.Loggers {
		buf, err := json.Marshal(v)
		checkErr(err)
		ret_p, err := conn.Create("/ymb/loggers/"+v.Id, buf, 0, zk.WorldACL(zk.PermAll))
		checkErr(err)
		fmt.Printf("ZNode[%s] created.\n", ret_p)
	}
}

// delete znode recur
func deleteRecur(conn *zk.Conn, path string) error {
	flag, _, err := conn.Exists(path)
	checkErr(err)
	if !flag {
		fmt.Print("ZNode deleted already.\n")
		return nil
	}
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err)
	fmt.Printf("Find %s\n", path)
	for _, znode := range children {
		sub_znode := path + "/" + znode
		deleteRecur(conn, sub_znode)
	}
	checkErr(conn.Delete(path, -1))
	return nil
}

// check and process errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
