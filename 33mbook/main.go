package main

import (
	"github.com/astaxie/beego"
	_ "study/golang/33mbook/routers"
	_ "study/golang/33mbook/sysint"
)

func main() {
	beego.Run()
}
