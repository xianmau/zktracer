package main

import (
	"database/sql"
	"datastructure"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
	"zk"
)

const (
	MYSQL_CONN_STR = "root:root@/ymb?charset=utf8"
)

var (
	zones = map[string][]string{
	//"zone1": []string{"172.19.32.16"},
	//"zone2": []string{"172.19.32.153"},
	//"zone3": []string{"172.19.32.46"},
	}
)

func main() {
	fmt.Printf("Data extractor started.\n")
	timer := time.Tick(60 * time.Second)
	for _ = range timer {
		// 先获取zones
		err := getZonesMap()
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		var wg sync.WaitGroup
		for zoneid, ips := range zones {
			wg.Add(1)
			go func(zoneid string, ips []string) {
				err := extracting(ips)
				if err != nil {
					log.Printf("%v\n", err)
				}
				log.Printf("zone [%s] done\n", zoneid)
				wg.Done()
			}(zoneid, ips)
		}
		wg.Wait()
		//runtime.GC()
		log.Printf("all zones done\n")
	}
}

// 从数据库中获取所有的zones
func getZonesMap() error {
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("select * from `tb_zone`")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var Id string
		var Ip string
		err := rows.Scan(&Id, &Ip)
		if err != nil {
			return err
		}
		var ipSet []string
		err = json.Unmarshal([]byte(Ip), &ipSet)
		if err != nil {
			return err
		}
		zones[Id] = ipSet
	}
	return nil
}

func extracting(host []string) error {
	conn := zk.Connect(host, time.Second)
	defer conn.Close()

	// 实例化一个ymb，把从zookeeper上的数据先放到里面
	t1 := time.Now()
	ymb := datastructure.YMB{
		[]datastructure.BrokerInfo{},
		[]datastructure.LoggerInfo{},
		[]datastructure.AppInfo{},
		[]datastructure.TopicInfo{},
		string(""),
		[]string{},
		[]datastructure.BrokerTopic{},
	}
	extractingZoneId(conn, &ymb)
	extractingRemoteZones(conn, &ymb)
	extractingApps(conn, &ymb)
	extractingBrokers(conn, &ymb)
	extractingLoggers(conn, &ymb)
	extractingTopics(conn, &ymb)
	t2 := time.Now()
	log.Printf("Extract data from [%s] using %v\n", ymb.ZoneId, t2.Sub(t1))

	// 再把ymb持久化到数据库
	t1 = time.Now()
	persistToLocalStorageUsingTx(&ymb)
	t2 = time.Now()
	log.Printf("Store data in database using %v\n", t2.Sub(t1))

	return nil
}

// extract zoneid znodes' data
func extractingZoneId(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/zoneid"
	if flag, err := conn.Exists(path); err == nil && flag {
		if data, err := conn.Get(path); err == nil {
			ymb.ZoneId = string(data)
		}
	}
}

// extract remote_zones znodes' data
func extractingRemoteZones(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/remote_zones"
	if flag, err := conn.Exists(path); err == nil && flag {
		if data, err := conn.Get(path); err == nil {
			var item []string
			json.Unmarshal([]byte(data), &item)
			ymb.RemoteZones = item
		}
	}
}

// extract app znodes' data
func extractingApps(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/appid"
	children, err := conn.Children(path)
	if err != nil {
		log.Printf("zone [%s]: %v\n", err)
	}

	for _, znode := range children {
		cur_znode := path + "/" + znode
		if data, err := conn.Get(cur_znode); err == nil {
			var item datastructure.AppInfo
			json.Unmarshal([]byte(data), &item)
			ymb.Apps = append(ymb.Apps, item)
		}
	}
}

// extract logger znodes' data
func extractingLoggers(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/loggers"
	children, err := conn.Children(path)
	if err != nil {
		log.Printf("zone [%s]: %v\n", err)
	}

	for _, znode := range children {
		cur_znode := path + "/" + znode
		if data, err := conn.Get(cur_znode); err == nil {
			var item datastructure.LoggerInfo
			json.Unmarshal([]byte(data), &item)
			ymb.Loggers = append(ymb.Loggers, item)
		}
	}
}

// extract broker znodes' data
func extractingBrokers(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/brokers"
	children, err := conn.Children(path)
	if err != nil {
		log.Printf("zone [%s]: %v\n", err)
	}

	for _, znode := range children {
		cur_znode := path + "/" + znode
		if data, err := conn.Get(cur_znode); err == nil {
			var item datastructure.BrokerInfo
			json.Unmarshal([]byte(data), &item)
			ymb.Brokers = append(ymb.Brokers, item)
		}
	}
}

// extract topic znodes' data
func extractingTopics(conn *zk.ZK, ymb *datastructure.YMB) {
	path := "/ymb/topics"
	children, err := conn.Children(path)
	if err != nil {
		log.Printf("zone [%s]: %v\n", err)
	}

	for _, znode := range children {
		sub_path := path + "/" + znode
		children2, err := conn.Children(sub_path)
		if err != nil {
			log.Printf("zone [%s]: %v\n", err)
		}

		for _, sub_znode := range children2 {
			sub_sub_path := sub_path + "/" + sub_znode
			if data, err := conn.Get(sub_sub_path); err == nil {
				var item datastructure.TopicInfo
				json.Unmarshal([]byte(data), &item)
				ymb.Topics = append(ymb.Topics, item)

				children3, err := conn.Children(sub_sub_path)
				if err != nil {
					log.Printf("zone [%s]: %v\n", err)
				}
				for _, sub_sub_znode := range children3 {
					var item2 = datastructure.BrokerTopic{}
					item2.AppId = znode
					item2.TopicName = sub_znode
					item2.BrokerId = sub_sub_znode
					ymb.BrokersTopics = append(ymb.BrokersTopics, item2)
				}
			}
		}
	}
}

// persist object to local storage
func persistToLocalStorageUsingTx(ymb *datastructure.YMB) {
	var (
		db   *sql.DB
		stmt *sql.Stmt
		tx   *sql.Tx
	)
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err, "Store:Connect error")

	tx, err = db.Begin()
	checkErr(err, "Store:Tx begin error")
	// store app
	_, err = tx.Exec("update `tb_app` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	checkErr(err, "Store:App")
	for _, e := range ymb.Apps {
		stmt, err = tx.Prepare("insert into `tb_app`(`Id`,`ZoneId`,`Key`,`Status`) values (?,?,?,?) on duplicate key update `Key`=values(`Key`),`Status`=true")
		checkErr(err, "")
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Key, true)
		checkErr(err, "")
	}
	// store broker
	_, err = tx.Exec("update `tb_broker` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	checkErr(err, "")
	for _, e := range ymb.Brokers {
		stmt, err = tx.Prepare("insert into `tb_broker`(`Id`,`ZoneId`,`Addrs`,`Status`) values (?,?,?,?) on duplicate key update `Addrs`=values(`Addrs`),`Status`=true")
		checkErr(err, "")
		addrs_json, _ := json.Marshal(e.Addrs)
		_, err = stmt.Exec(e.Id, ymb.ZoneId, string(addrs_json), true)
		checkErr(err, "")
	}

	// store broker_stat
	for _, e := range ymb.Brokers {
		stmt, err = tx.Prepare("insert into `tb_broker_stat`(`BrokerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		checkErr(err, "")
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		checkErr(err, "")
	}

	// store logger
	_, err = tx.Exec("update `tb_logger` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	checkErr(err, "")
	for _, e := range ymb.Loggers {
		stmt, err = tx.Prepare("insert into `tb_logger`(`Id`,`ZoneId`,`Addr`,`BlkDev`,`Status`) values (?,?,?,?,?) on duplicate key update `Addr`=values(`Addr`),`BlkDev`=values(`BlkDev`),`Status`=true")
		checkErr(err, "")
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Addr, e.BlkDev, true)
		checkErr(err, "")
	}

	// store logger_stat
	for _, e := range ymb.Loggers {
		stmt, err = tx.Prepare("insert into `tb_logger_stat`(`LoggerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		checkErr(err, "")
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		checkErr(err, "")
	}

	// store topic
	_, err = tx.Exec("update `tb_topic` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	checkErr(err, "")
	for _, e := range ymb.Topics {
		stmt, err = tx.Prepare("insert into `tb_topic`(`Name`,`AppId`,`ZoneId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status`) values (?,?,?,?,?,?,?,?) on duplicate key update `BrokerId`=values(`BrokerId`),`ReplicaNum`=values(`ReplicaNum`),`Retention`=values(`Retention`),`Segments`=values(`Segments`),`Status`=true")
		checkErr(err, "")
		segments_json, _ := json.Marshal(e.Segments)
		_, err = stmt.Exec(e.Name, e.AppId, ymb.ZoneId, GetBrokerId(ymb, e.Name, e.AppId), e.ReplicaNum, e.ReplicaNum, segments_json, true)
		checkErr(err, "")
	}
	err = tx.Commit()
	checkErr(err, "")

	db.Close()
}

func GetBrokerId(ymb *datastructure.YMB, Name string, AppId string) string {
	for _, v := range ymb.BrokersTopics {
		if v.TopicName == Name && v.AppId == AppId {
			return v.BrokerId
		}
	}
	return ""
}

// traverse all znodes under the specified path
func traverse(conn *zk.ZK, path string) {
	children, err := conn.Children(path)
	checkErr(err, "")

	if len(children) <= 0 {
		data, err := conn.Get(path)
		if err == nil {
			fmt.Println("#Leaf ZNode Found:")
			fmt.Println("#PATH: ", path)
			fmt.Println("#DATA: ", string(data))
		}
	}
	for _, znode := range children {
		if path == "/" {
			fmt.Printf("Searching ZNode: /%s\n", znode)
			traverse(conn, "/"+znode)
		} else {
			fmt.Printf("Searching ZNode: %s/%s\n", path, znode)
			traverse(conn, path+"/"+znode)
		}
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Printf("%v", err)
	}
}
