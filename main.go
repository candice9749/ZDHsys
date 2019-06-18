package main

import (
	_ "ZDHsys/models"
	_ "ZDHsys/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

