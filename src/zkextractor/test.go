package main

import (
	"database/sql"
	"datastructure"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samuel/go-zookeeper/zk"
	"runtime"
	"time"
)

const (
	MYSQL_CONN_STR = "root:root@/ymb?charset=utf8"
	FIRST_IP       = "172.19.32.153"
)

var (
	ymb         datastructure.YMB
	brokertopic []datastructure.BrokerTopic
)

var (
	HOSTS = [][]string{
		{"172.19.32.16", "172.19.32.153", "172.19.32.46"},
		{"172.19.32.35", "172.19.32.135", "172.19.32.124"},
	}
)

func main() {
	/*-------------------------------------------------------------------------------------
	 * 流程说明：
	 * 1，通过第一个zone获取所有zone的信息，存到zones数组中，并存至数据库
	 * 2，设置定时器，周期暂时定为1分钟，定时执行
	 *    1）从数据库读取zones数据到zones数组中，用于当前遍历所有zone的依据
	 *    2）遍历zones，对每个zone
	 *       a）提取zoneid和remote_zones，并对数据库中的zone信息进行修正
	 *       b）提取app，broker，logger，topic，并存至数据库中
	 *------------------------------------------------------------------------------------*/

	fmt.Printf("Extractor Started.\n\n")
	// get all zones' data by the FIRST_IP
	GetAllZonesFromZK([]string{FIRST_IP})

	// extract data per minute
	timer := time.Tick(2 * time.Second)
	for _ = range timer {
		// traverse all zones
		runtime.GC()
		go func() {
			zones := GetAllZonesFromDB()
			fmt.Printf("Current zones: %s\n\n", zones)
			for i, ip := range zones {
				InitYMB()
				fmt.Printf("Extracting zone[%s] at %s\n", ip, time.Now().Format("2006-01-02 15:04:05"))
				fmt.Printf("--------------------------------------------------------------------------\n")
				if extracting(HOSTS[i]) {
					fmt.Printf("----\ndone\n\n")
				} else {
					fmt.Printf("-----\nfaild\n\n")
				}
			}
		}()

	}
}

func InitYMB() {
	ymb = datastructure.YMB{
		[]datastructure.BrokerInfo{},
		[]datastructure.LoggerInfo{},
		[]datastructure.AppInfo{},
		[]datastructure.TopicInfo{},
		string(""),
		[]string{},
	}
	brokertopic = []datastructure.BrokerTopic{}
}

func GetAllZonesFromZK(host []string) {
	zones := []string{}
	conn, _, err := zk.Connect(host, time.Second)
	checkErr(err, "GetAllZonesFromZK")
	if flag, _, err := conn.Exists("/ymb/zoneid"); err == nil && flag {
		if data, _, err := conn.Get("/ymb/zoneid"); err == nil {
			var item string
			json.Unmarshal([]byte(data), &item)
			zones = append(zones, item)
		}
	}
	if flag, _, err := conn.Exists("/ymb/remote_zones"); err == nil && flag {
		if data, _, err := conn.Get("/ymb/remote_zones"); err == nil {
			var item []string
			json.Unmarshal([]byte(data), &item)
			for _, v := range item {
				zones = append(zones, v)
			}
		}
	}
	// storage data into mysql
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err, "GetAllZonesFromZK:Open mysql error.")
	_, err = db.Exec("delete from `tb_zone` where true")
	checkErr(err, "GetAllZonesFromZK:Execute sql error.")
	for _, v := range zones {
		_, err = db.Exec("insert into `tb_zone`(`Id`) values (?) on duplicate key update `Id`=values(`Id`)", v)
		checkErr(err, "GetAllZonesFromZK:Insert node error.")
	}
	db.Close()
}

func GetAllZonesFromDB() []string {
	zones := []string{}
	// storage data into mysql
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err, "GetAllZonesFromDB")
	rows, err := db.Query("select * from `tb_zone`")
	checkErr(err, "GetAllZonesFromDB")
	defer rows.Close()
	for rows.Next() {
		var Id string
		checkErr(rows.Scan(&Id), "")
		zones = append(zones, Id)
	}
	db.Close()
	return zones
}

func extracting(host []string) bool {
	//when host can't be connected, return false
	conn, _, err := zk.Connect(host, time.Second)
	checkErr(err, "")

	fmt.Printf("DEBUG[EXTRACT-START]%s\n", time.Now().Format("2006-01-02 15:04:05"))

	extractingZoneId(conn, "/ymb/zoneid")
	fmt.Printf("1")
	extractingRemoteZones(conn, "/ymb/remote_zones")
	fmt.Printf("2")
	fixZones()
	fmt.Printf("3")

	extractingApps(conn, "/ymb/appid")
	fmt.Printf("4")
	extractingBrokers(conn, "/ymb/brokers")
	fmt.Printf("5")
	extractingLoggers(conn, "/ymb/loggers")
	fmt.Printf("6")
	extractingTopics(conn, "/ymb/topics")

	fmt.Printf("DEBUG[EXTRACT-END  ]%s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	//persistToLocalStorageUsingTx()
	return true
}

// extract zoneid znodes' data
func extractingZoneId(conn *zk.Conn, path string) {
	if flag, _, err := conn.Exists(path); err == nil && flag {
		if data, _, err := conn.Get(path); err == nil {
			var item string
			json.Unmarshal([]byte(data), &item)
			ymb.ZoneId = item
		}
	}
}

// extract remote_zones znodes' data
func extractingRemoteZones(conn *zk.Conn, path string) {
	if flag, _, err := conn.Exists(path); err == nil && flag {
		if data, _, err := conn.Get(path); err == nil {
			var item []string
			json.Unmarshal([]byte(data), &item)
			ymb.RemoteZones = item
		}
	}
}

// extract app znodes' data
func extractingApps(conn *zk.Conn, path string) {
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err, "exApp")

	for _, znode := range children {
		cur_znode := path + "/" + znode
		data, _, err := conn.Get(cur_znode)
		if err != nil {
			fmt.Printf("Get %s's data faild.", cur_znode)
			continue
		}
		var item datastructure.AppInfo
		json.Unmarshal([]byte(data), &item)
		ymb.Apps = append(ymb.Apps, item)
	}
}

// extract logger znodes' data
func extractingLoggers(conn *zk.Conn, path string) {
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err, "exLogger")

	for _, znode := range children {
		cur_znode := path + "/" + znode
		data, _, err := conn.Get(cur_znode)
		if err != nil {
			fmt.Printf("Get %s's data faild.", cur_znode)
			continue
		}
		var item datastructure.LoggerInfo
		json.Unmarshal([]byte(data), &item)
		ymb.Loggers = append(ymb.Loggers, item)
	}
}

// extract broker znodes' data
func extractingBrokers(conn *zk.Conn, path string) {
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err, "exBroker")

	for _, znode := range children {
		cur_znode := path + "/" + znode
		data, _, err := conn.Get(cur_znode)
		if err != nil {
			fmt.Printf("Get %s's data faild.", cur_znode)
			continue
		}
		var item datastructure.BrokerInfo
		json.Unmarshal([]byte(data), &item)
		ymb.Brokers = append(ymb.Brokers, item)
	}
}

// extract topic znodes' data
func extractingTopics(conn *zk.Conn, path string) {
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err, "exTopic")

	for _, znode := range children {
		sub_path := path + "/" + znode
		children2, _, _, err := conn.ChildrenW(sub_path)
		checkErr(err, "exTopic")
		// do something

		for _, sub_znode := range children2 {
			sub_sub_path := sub_path + "/" + sub_znode
			data, _, err := conn.Get(sub_sub_path)
			if err != nil {
				fmt.Printf("Get $'s data faild.\n", sub_sub_path)
				continue
			}
			var item datastructure.TopicInfo
			json.Unmarshal([]byte(data), &item)
			ymb.Topics = append(ymb.Topics, item)

			children3, _, _, err := conn.ChildrenW(sub_sub_path)
			checkErr(err, "exTopic")
			for _, sub_sub_znode := range children3 {
				var item2 = datastructure.BrokerTopic{}
				item2.AppId = znode
				item2.TopicName = sub_znode
				item2.BrokerId = sub_sub_znode
				brokertopic = append(brokertopic, item2)
			}
		}
	}
}

func fixZones() {
	var (
		db *sql.DB
		//stmt *sql.Stmt
		//res  sql.Result
	)
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err, "")
	_, err = db.Exec("insert into `tb_zone`(`Id`) values (?) on duplicate key update `Id`=values(`Id`)", ymb.ZoneId)
	checkErr(err, "fixZones")

	for _, e := range ymb.RemoteZones {
		_, err = db.Exec("insert into `tb_zone`(`Id`) values (?) on duplicate key update `Id`=values(`Id`)", e)
		checkErr(err, "fixZones")
	}
	db.Close()
}

// persist object to local storage
func persistToLocalStorageUsingTx() {
	var (
		db   *sql.DB
		stmt *sql.Stmt
		tx   *sql.Tx
	)
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err, "Store:Connect error")

	fmt.Printf("DEBUG[STORE-START]%s\n", time.Now().Format("2006-01-02 15:04:05"))
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
		_, err = stmt.Exec(e.Name, e.AppId, ymb.ZoneId, GetBrokerId(e.Name, e.AppId), e.ReplicaNum, e.ReplicaNum, segments_json, true)
		checkErr(err, "")
	}
	err = tx.Commit()
	checkErr(err, "")
	fmt.Printf("DEBUG[STORE-END  ]%s\n", time.Now().Format("2006-01-02 15:04:05"))
	db.Close()
}

func GetBrokerId(Name string, AppId string) string {
	for _, v := range brokertopic {
		if v.TopicName == Name && v.AppId == AppId {
			return v.BrokerId
		}
	}
	return ""
}

// traverse all znodes under the specified path
func traverse(conn *zk.Conn, path string) {
	children, stat, _, err := conn.ChildrenW(path)

	checkErr(err, "")

	if stat.NumChildren <= 0 {
		data, _, err := conn.Get(path)
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

// check and process errors
func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf("DEBUG: %s_____________________________", msg)
		panic(err)
	}
}
