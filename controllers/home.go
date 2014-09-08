package controllers

import ()

type HomeCtrl struct {
	baseController
}

func (this *HomeCtrl) Get() {
	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	this.TplNames = "index.tpl"

}
