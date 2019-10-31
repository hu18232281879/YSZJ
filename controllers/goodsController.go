package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/models"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndex() {
	userName := this.GetSession("pyg_userName")
	if userName != nil {
		this.Data["userName"] = userName
	}
	var goodsTypes []map[string]interface{}
	o := orm.NewOrm()
	var goodsMenu []models.TpshopCategory
	o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&goodsMenu)
	for _, yiji := range goodsMenu {
		tempContainer := make(map[string]interface{})
		var erji []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid", yiji.Id).All(&erji)
		tempContainer["yiji"] = yiji
		tempContainer["erji"] = erji

		goodsTypes = append(goodsTypes, tempContainer)
	}
	for _, v := range goodsTypes {
		var erjiContainer []map[string]interface{}
		for _, erClass := range v["erji"].([]models.TpshopCategory) {
			var sanji []models.TpshopCategory
			tempContainer := make(map[string]interface{})
			o.QueryTable("TpshopCategory").Filter("Pid", erClass.Id).All(&sanji)
			tempContainer["erji"] = erClass
			tempContainer["sanji"] = sanji
			//把三级容器放到二级容器
			erjiContainer = append(erjiContainer, tempContainer)
		}
		v["sanji"] = erjiContainer
	}

	this.Data["goodsTypes"] = goodsTypes

	this.TplName = "index.html"
}
func (this *GoodsController) ShowIndexSx() {
	this.TplName="index_sx.html"
}
