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

	LoginName := this.GetSession("admin")
	this.Data["IsLogin"] = LoginName != nil
	this.Data["LoginName"] = LoginName
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

func (this *AdminController) GetNodeData() {
	currentnode := this.GetString("znode")
	currentzone := this.GetString("zoneid")
	conn, _, err := zk.Connect([]string{currentzone}, time.Second)
	defer conn.Close()
	checkErr(err)
	data, stat, _, err := conn.GetW(currentnode)
	checkErr(err)
	if data == nil || len(data) <= 0 || data[0] == 0 {
		this.Ctx.WriteString("[]")
	}
	fmt.Printf("Stat -> %-v\n", stat)
	this.Ctx.WriteString(string(data))
}

func (this *AdminController) CreateNode() {
	nodepath := this.Input().Get("nodepath")
	newnode := this.Input().Get("znode")
	nodedata := this.Input().Get("data")
	currentzone := this.Input().Get("zoneid")

	data := []zktreenode{}
	conn, _, err := zk.Connect([]string{currentzone}, time.Second)
	defer conn.Close()
	if nodepath == "/ymb" {
		this.Ctx.WriteString("[]")
		return
	}
	fmt.Println("the data: ", nodedata)
	path, err := conn.Create(nodepath+"/"+newnode, []byte(nodedata), 0, zk.WorldACL(zk.PermAll))
	checkErr(err)
	zkn := zktreenode{newnode, true}
	data = append(data, zkn)
	datajson, err := json.Marshal(data)
	this.Ctx.WriteString(string(datajson))
	fmt.Println(path)
}

func (this *AdminController) DeleteNode() {
	node := this.Input().Get("node")
	zone := this.Input().Get("zoneid")
	fmt.Println(node, zone)
	conn, _, err := zk.Connect([]string{zone}, time.Second)
	defer conn.Close()
	if node == "/ymb" {
		this.Ctx.WriteString("You shouldn't delete the root!!!")
		return
	}
	err = deleteRecur(conn, node)
	checkErr(err)
	this.Ctx.WriteString("Deleted.")
}

func deleteRecur(conn *zk.Conn, path string) error {
	flag, _, err := conn.Exists(path)
	checkErr(err)
	if !flag {
		fmt.Print("ZNode deleted already.\n")
		return nil
	}
	children, _, _, err := conn.ChildrenW(path)
	checkErr(err)
	//fmt.Printf("Find %s\n", path)
	for _, znode := range children {
		sub_znode := path + "/" + znode
		deleteRecur(conn, sub_znode)
	}
	checkErr(conn.Delete(path, -1))
	return nil
}

func (this *AdminController) Zone() {
	this.Layout = "layout.tpl"
	this.TplNames = "admin/zone.html"

	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "home/styles.tpl"
	this.LayoutSections["Scripts"] = "admin/zonescripts.tpl"

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

	this.Data["zonelist"] = zonelist
	rzjson, err := json.Marshal(zonelist)
	checkErr(err)
	this.Data["remotezones"] = string(rzjson)
	LoginName := this.GetSession("admin")
	this.Data["IsLogin"] = LoginName != nil
	this.Data["LoginName"] = LoginName

}

func (this *AdminController) CreateZone() {
	zoneid := this.Input().Get("newzoneid")
	remotezones := this.Input().Get("remotezones")
	conn, _, err := zk.Connect([]string{zoneid}, time.Second)
	checkErr(err)
	_, err = conn.Create("/ymb/zoneid", []byte(zoneid), 0, zk.WorldACL(zk.PermAll))
	checkErr(err)
	_, err = conn.Create("/ymb/remote_zones", []byte(remotezones), 0, zk.WorldACL(zk.PermAll))
	checkErr(err)
	conn.Close()

	var rz []string
	json.Unmarshal([]byte(remotezones), &rz)

	for _, z := range rz {
		conn, _, err := zk.Connect([]string{z}, time.Second)
		checkErr(err)
		data, _, err := conn.Get("/ymb/remote_zones")
		var newrz []string
		json.Unmarshal(data, &newrz)
		newrz = append(newrz, zoneid)
		buf, err := json.Marshal(newrz)
		checkErr(err)
		_, err = conn.Set("/ymb/remote_zones", buf, 0)
		checkErr(err)
	}
}

func (this *AdminController) DeleteZone() {
	zoneid := this.Input().Get("zoneid")
	remotezones := this.Input().Get("remotezones")
	var rz []string
	json.Unmarshal([]byte(remotezones), &rz)

	for _, z := range rz {
		if z == zoneid {
			continue
		}
		conn, _, err := zk.Connect([]string{z}, time.Second)
		checkErr(err)
		var newrz []string
		json.Unmarshal([]byte(remotezones), &newrz)

		index := -1
		for k, v := range newrz {
			if v == zoneid {
				index = k
				break
			}
		}
		if index == -1 {
			return
		}
		newrz = append(newrz[:index], newrz[index+1:]...)
		buf, err := json.Marshal(newrz)
		checkErr(err)
		_, err = conn.Set("/ymb/remote_zones", buf, 0)
		checkErr(err)
	}

}
