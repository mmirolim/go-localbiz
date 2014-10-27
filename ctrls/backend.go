package ctrls

type Backend struct {
	baseController
}

func (c *Backend) DashBoard() {
	c.TplNames = "backend/dashboard.tpl"
	c.Data["Title"] = "Welcome to Dashboard"
}
