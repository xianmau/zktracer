package controllers

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//"math/rand"
	"time"
	"zktracer/models"
)

type AppController struct {
	beego.Controller
}

// 显示指定zone下的所有app的最新状态
func (this *AppController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "app/app.html"
	this.LayoutSections = make(map[string]string)
	//this.LayoutSessions["Styles"] = "app/styles.tpl"
	this.LayoutSections["Scripts"] = "app/scripts.tpl"

	applist := []models.AppModel{}
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

	rows, _ = db.Query("select `Id`,`Key`,`Status` from `tb_app` where `ZoneId`=?", currentzone)
	for rows.Next() {
		tmp := models.AppModel{}
		checkErr(rows.Scan(&tmp.Id, &tmp.Key, &tmp.Status))
		applist = append(applist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["currentzone"] = currentzone
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["zonelist"] = zonelist
	this.Data["applist"] = applist

	fmt.Println("a client request.")
}
