package main

import (
	"github.com/astaxie/beego"
	_ "zktracer/routers"
)

func main() {
	beego.Run()
}
