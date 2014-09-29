package ctrls

type Home struct {
	baseController
}

func (this *Home) Get() {
	this.TplNames = "index.tpl"

	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	isAuth := this.GetSession("isAuth")
	if isAuth != nil {
		this.Data["isAuth"] = isAuth.(bool)
	}
}
