package sysint

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func dbinit(alias string) {
	dbAlias := alias

	if "w" == alias || "default" == alias || len(alias) < 0 {
		dbAlias = "default"
		alias = "w"
	}

	dbName := beego.AppConfig.String("db_" + alias + "_database")
	dbUser := beego.AppConfig.String("db_" + alias + "_username")
	dbPwd := beego.AppConfig.String("db_" + alias + "_password")
	dbHost := beego.AppConfig.String("db_" + alias + "_host")
	dbPort := beego.AppConfig.String("db_" + alias + "_port")
	orm.RegisterDataBase(dbAlias, "mysql", dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8", 30)

	isDev := "dev" == beego.AppConfig.String("runmode")

	if "w" == alias {
		orm.RunSyncdb("default", false, false)
	}

	if isDev {
		orm.Debug = isDev
	}
}
