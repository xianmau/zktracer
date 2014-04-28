package controllers

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Layout = "layout.tpl"
	this.TplNames = "home/home.html"

	this.LayoutSections = make(map[string]string)
	//this.LayoutSections["Styles"] = "home/styles.tpl"
	this.LayoutSections["Scripts"] = "home/scripts.tpl"
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
