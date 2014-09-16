package ctrls

import (
//"github.com/astaxie/beego"
)

type Home struct {
	baseController
}

func (this *Home) Get() {
	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	this.Data["Num"] = 10
	this.SetSession("num", 12)
	this.TplNames = "index.tpl"

}
