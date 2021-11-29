package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/i18n"

)
type AccountResponse struct{
	Id int64
	Name    string
	Email   string
}
type AccountController struct {
	BaseController
	i18n.Locale
}
func (a *AccountController) Prepare() {
    a.EnableXSRF = false
}
///valid account
func (c *AccountController) Validaccount() {
	username := c.GetString("username")
	pass:=c.GetString("pass")
	if len(username) == 0 || len(pass) == 0{
		//用户名和邮箱为空
		c.ErrorJson(20211122163020,c.Tr("email_pass_empty"),nil)
	}
	
	
	account,err:=models.Validaccount(username,pass)

	if err !=nil {	
		
		c.ErrorJson(202108031641149,err.Error(),nil)
	} else {
		c.SetSession("uid", account.Id)
		models.DefaultAccountLoginLog.AccountLogin(&account)
		accountRes :=AccountResponse{Id:account.Id,
			Name:account.Name,
			Email: account.Email,
		}

		c.SuccessJson(accountRes)
	}
}

