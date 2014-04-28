package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "admin/admin.html"

	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "home/styles.tpl"
	this.LayoutSections["Scripts"] = "admin/scripts.tpl"

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

	db.Close()

	this.Data["currentzone"] = currentzone
	this.Data["zonelist"] = zonelist
}

type zktreenode struct {
	Text     string `json:"text"`
	Children bool   `json:"children"`
}

func (this *AdminController) GetData() {

	currentnode := "/" + this.GetString("znode")
	currentzone := this.GetString("zoneid")
	fmt.Println(currentnode)
	data := []zktreenode{}
	conn, _, err := zk.Connect([]string{currentzone}, time.Second)
	checkErr(err)
	if flag, _, err := conn.Exists(currentnode); err == nil && flag {
		children, _, _, err := conn.ChildrenW(currentnode)
		checkErr(err)
		for _, znode := range children {
			fmt.Println(znode)
			d := zktreenode{znode, true}
			data = append(data, d)
		}
	}
	fmt.Println(data)
	datajson, err := json.Marshal(data)
	checkErr(err)
	fmt.Println(datajson)

	this.Ctx.WriteString(string(datajson))
}
