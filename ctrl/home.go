package ctrl

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
	this.TplNames = "index.tpl"

}
