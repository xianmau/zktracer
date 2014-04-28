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

type LoggerController struct {
	beego.Controller
}

func (this *LoggerController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "logger/logger.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSession["Styles"] = "logger/styles.tpl"
	this.LayoutSections["Scripts"] = "logger/scripts.tpl"

	loggerlist := []models.LoggerModel{}
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

	rows, _ = db.Query("select A.`Id`,A.`Addr`,A.`BlkDev`,A.`Status`,B.`Cpu`,B.`Net`,B.`Disk` from `tb_logger` as A left join `tb_logger_stat` as B on A.`Id`=B.`LoggerId` and A.`ZoneId`=B.`ZoneId` and (B.`Timestamp`>=now() - interval 1 minute) where A.`ZoneId`=?", currentzone)
	for rows.Next() {
		tmp := models.LoggerModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		checkErr(rows.Scan(&tmp.Id, &tmp.Addr, &tmp.BlkDev, &tmp.Status, &Cpu, &Net, &Disk))

		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		loggerlist = append(loggerlist, tmp)

	}
	rows.Close()
	db.Close()
	this.Data["currentzone"] = currentzone
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["zonelist"] = zonelist
	this.Data["loggerlist"] = loggerlist

	fmt.Println("a client request.")
}

func (this *LoggerController) Detail() {
	this.Layout = "layout.tpl"
	this.TplNames = "logger/detail.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSessions["Styles"] = "logger/styles.tpl"
	this.LayoutSections["Scripts"] = "logger/scripts.tpl"

	currentzone := this.GetString("zoneid")
	currentlogger := this.GetString("loggerid")
	statdata := []models.LoggerStatModel{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_logger_stat` where `LoggerId`=? and `ZoneId`=? order by `Timestamp` asc", currentlogger, currentzone)
	for rows.Next() {
		tmp := models.LoggerStatModel{}
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
	this.Data["currentlogger"] = currentlogger
	this.Data["statdata"] = statdata

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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

// 获取指定zone和logger的最新数据，用于实时数据展示
func (this *LoggerController) GetLatestData() {

	zoneid := this.GetString("zoneid")
	loggerid := this.GetString("loggerid")
	// 创建随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_logger_stat` where `LoggerId`=? and `ZoneId`=? and (`Timestamp`>=now() - interval 1 minute)", loggerid, zoneid)
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
