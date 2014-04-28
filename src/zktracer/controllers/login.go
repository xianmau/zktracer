package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "login.html"

	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "home/styles.tpl"
	//this.LayoutSections["Scripts"] = "home/scripts.tpl"
}

func (this *LoginController) Post() {
	name := this.Input().Get("Name")
	password := this.Input().Get("Password")
	//remenber := this.Input.Get("Remenber")
	if name == "admin" && password == "admin" {

		this.SetSession("admin", string(name))

		recvName := this.GetSession("admin")
		fmt.Printf("%s", recvName)

		this.Ctx.Redirect(302, "/admin")

		//this.Ctx.WriteString("success, " + recvName.(string))
	} else {
		this.Ctx.WriteString("faild")
	}
}
