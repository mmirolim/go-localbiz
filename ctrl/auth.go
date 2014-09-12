package ctrl

type Auth struct {
	baseController
}

func (this *Auth) Get() {
	this.TplNames = "login.tpl"
}
