package controllers

import (
	"github.com/astaxie/beego"
	"study/golang/33mbook/models"
)

type ExploreController struct {
	BaseController
}

func (c *ExploreController) Index() {
	var (
		cid       int
		cate      models.Category
		urlPrefix = beego.URLFor("ExploreController.Index")
	)

	if cid, _ = c.GetInt("cid"); cid > 0 {
		cateModel := new(models.Category)
		cate = cateModel.Find(cid)
		c.Data["Cate"] = cate
	}

	c.Data["Cid"] = cid
	c.TplName = "explore/index.html"

	pageIndex,_ := c.GetInt("page",1)

}
