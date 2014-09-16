package ctrls

type User struct {
	baseController
}

func (this *User) Get() {
	this.TplNames = "user/signup.tpl"
}

func (this *User) SignUp() {
	this.TplNames = "user/signup.tpl"
}
