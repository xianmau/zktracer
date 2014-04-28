package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strconv"
	"time"
	"zktracer/models"
)

type BrokerController struct {
	beego.Controller
}

// 显示指定zone下的所有broker的最新动态
func (this *BrokerController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "broker/broker.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "broker/styles.tpl"
	this.LayoutSections["Scripts"] = "broker/scripts.tpl"

	brokerlist := []models.BrokerModel{}
	zonelist := []string{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))

	rows, _ := db.Query("select * from `tb_zone`")
	for rows.Next() {
		var Id string
		checkErr(rows.Scan(&Id))
		zonelist = append(zonelist, Id)
	}
	rows.Close()
	currentzone := this.GetString("zoneid")
	if currentzone == "" {
		currentzone = zonelist[0]
	}

	rows, _ = db.Query("select A.`Id`,A.`Addrs`,A.`Status`,B.`Cpu`,B.`Net`,B.`Disk` from `tb_broker` as A left join `tb_broker_stat` as B on A.`Id`=B.`BrokerId` and A.`ZoneId`=B.`ZoneId` and (B.`Timestamp`>=now() - interval 1 minute) where A.`ZoneId`=?", currentzone)
	for rows.Next() {
		tmp := models.BrokerModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		checkErr(rows.Scan(&tmp.Id, &tmp.Addrs, &tmp.Status, &Cpu, &Net, &Disk))

		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		brokerlist = append(brokerlist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["currentzone"] = currentzone
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["zonelist"] = zonelist
	this.Data["brokerlist"] = brokerlist

	fmt.Println("a client request.")
}

// 显示某一个broker的详细性能数据
func (this *BrokerController) Detail() {
	this.Layout = "layout.tpl"
	this.TplNames = "broker/detail.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "broker/styles.tpl"
	this.LayoutSections["Scripts"] = "broker/scripts.tpl"

	currentzone := this.GetString("zoneid")
	currentbroker := this.GetString("brokerid")
	statdata := []models.BrokerStatModel{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_broker_stat` where `BrokerId`=? and `ZoneId`=? order by `Timestamp` asc", currentbroker, currentzone)
	for rows.Next() {
		tmp := models.BrokerStatModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		checkErr(rows.Scan(&tmp.Timestamp, &Cpu, &Net, &Disk))
		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		statdata = append(statdata, tmp)
	}
	db.Close()

	this.Data["currentzone"] = currentzone
	this.Data["currentbroker"] = currentbroker
	this.Data["statdata"] = statdata

	// 将数据处理成json格式可用于直接渲染曲线图
	// 创建随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 整CPU、NET、DISK数据
	cpuData := [][]int64{}
	netData := [][]int64{}
	diskData := [][]int64{}
	for _, v := range statdata {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy := []int64{timestamp.Unix() * 1000, int64((int(v.Cpu*100) + r.Intn(100)) % 100)}
		cpuData = append(cpuData, xy)

		timestamp, _ = time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy = []int64{timestamp.Unix() * 1000, int64((int(v.Net*100) + r.Intn(100)) % 100)}
		netData = append(netData, xy)

		timestamp, _ = time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy = []int64{timestamp.Unix() * 1000, int64((int(v.Disk*100) + r.Intn(100)) % 100)}
		diskData = append(diskData, xy)
	}
	tmp, _ := json.Marshal(cpuData)
	this.Data["cpuData"] = string(tmp)

	tmp, _ = json.Marshal(netData)
	this.Data["netData"] = string(tmp)

	tmp, _ = json.Marshal(diskData)
	this.Data["diskData"] = string(tmp)

	fmt.Println("a client request.")
}

// 获取指定zone和broker的最新数据，用于实时数据展示
func (this *BrokerController) GetLatestData() {

	zoneid := this.GetString("zoneid")
	brokerid := this.GetString("brokerid")
	// 创建随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_broker_stat` where `BrokerId`=? and `ZoneId`=? and (`Timestamp`>=now() - interval 1 minute)", brokerid, zoneid)
	for rows.Next() {
		var Timestamp string
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		checkErr(rows.Scan(&Timestamp, &Cpu, &Net, &Disk))
		timestamp, _ := time.Parse("2006-01-02 15:04:05", Timestamp)
		cpu, _ := strconv.ParseFloat(Cpu.String, 64)
		net, _ := strconv.ParseFloat(Net.String, 64)
		disk, _ := strconv.ParseFloat(Disk.String, 64)
		ret := "["
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64((int(cpu*100)+r.Intn(100))%100), 10) + "],"
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64((int(net*100)+r.Intn(100))%100), 10) + "],"
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64((int(disk*100)+r.Intn(100))%100), 10) + "]"
		ret += "]"
		this.Ctx.WriteString(ret)
	}
	this.Ctx.WriteString("")
	return
}
