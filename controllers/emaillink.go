package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
)

type EmailController struct {
	BaseController
	i18n.Locale
}

func (c *EmailController) Testreadsaveemail() {
	file:=""
	emailModel:=models.EmailLink{}
	lemail,lerr:=emailModel.ReademailFile(file)
	if(lerr!=nil){
		panic(lerr)
	}
	for _, le := range lemail {
		emailModel.SaveEmaildb(le)
	}
	c.SuccessJson(nil)
}