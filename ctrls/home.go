package ctrls

import (
)

type Home struct {
	baseController
}

func (this *Home) Get() {
	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	v := this.GetSession("number")
	if v == nil {
		this.SetSession("number", int(1))
		this.Data["Num"] = 0
	} else {
		this.SetSession("number", v.(int)+1)
		this.Data["Num"] = v.(int)
	}
	this.Data["Person"] = map[string]interface {}{"Person":"Mirolim"}
	this.TplNames = "index.tpl"

}
