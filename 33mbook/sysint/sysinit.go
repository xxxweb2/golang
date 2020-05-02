package sysint

import (
	"github.com/astaxie/beego"
	"path/filepath"
)

func sysinit() {
	uploads := filepath.Join("./", "uploads")
	beego.BConfig.WebConfig.StaticDir["/uploads"] = uploads

	//注册前端使用模块
	registerFunctions()
}

func registerFunctions() {

}
