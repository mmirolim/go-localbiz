package ctrls

type Auth struct {
	baseController
}

func (this *Auth) Get() {
	this.TplNames = "login.tpl"
}
