package controllers

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	"time"
	"zktracer/models"
)

type TopicController struct {
	beego.Controller
}

//
func (this *TopicController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "topic/topic.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "topic/style.tpl"
	this.LayoutSections["Scripts"] = "topic/scripts.tpl"

	topiclist := []models.TopicModel{}
	zonelist := []string{}
	applist := []string{}
	brokerlist := []string{}

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

	rows, _ = db.Query("select `Id` from `tb_app` where `ZoneId`=?", currentzone)
	for rows.Next() {
		var Id string
		checkErr(rows.Scan(&Id))
		applist = append(applist, Id)
	}
	rows.Close()

	rows, _ = db.Query("select `Id` from `tb_broker` where `ZoneId`=?", currentzone)
	for rows.Next() {
		var Id string
		checkErr(rows.Scan(&Id))
		brokerlist = append(brokerlist, Id)
	}
	rows.Close()

	currentapp := this.GetString("appid")
	currentbroker := this.GetString("brokerid")

	if currentapp != "" {
		rows, _ = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `ZoneId`=? and `AppId`=?", currentzone, currentapp)
	} else if currentbroker != "" {
		rows, _ = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `BrokerId`=?", currentbroker)
	} else {
		rows, _ = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `ZoneId`=?", currentzone)
	}
	for rows.Next() {
		tmp := models.TopicModel{}
		checkErr(rows.Scan(&tmp.Name, &tmp.AppId, &tmp.BrokerId, &tmp.ReplicaNum, &tmp.Retention, &tmp.Segments, &tmp.Status))
		topiclist = append(topiclist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["currentzone"] = currentzone
	this.Data["currentapp"] = currentapp
	this.Data["currentbroker"] = currentbroker
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["zonelist"] = zonelist
	this.Data["topiclist"] = topiclist
	this.Data["applist"] = applist
	this.Data["brokerlist"] = brokerlist

	fmt.Println("a client request.")
}
