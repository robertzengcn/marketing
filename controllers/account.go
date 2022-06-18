package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
)
type AccountResponse struct{
	Id int64
	Name    string
	Email   string
	Token   string
}
type AccountController struct {
	BaseController
	i18n.Locale
}
// func (a *AccountController) Prepare() {
//     a.EnableXSRF = false
// }
 func (c *AccountController) ChildPrepare(){
	// l := logs.GetLogger()
    //      l.Println("22222")
 }
///valid account
func (c *AccountController) Validaccount() {
	username := c.GetString("username")
	pass:=c.GetString("pass")
	if len(username) == 0 || len(pass) == 0{
		//用户名和邮箱为空
		c.ErrorJson(20211122163020,c.Tr("email_pass_empty"),nil)
	}
	// l := logs.GetLogger()
   
	// l.Println("33333")	
	// l.Println(c.Tr("welcome"))
	// l.Println("44444")
	account,err:=models.Validaccount(username,pass)

	if err !=nil {	
		
		c.ErrorJson(202108031641149,err.Error(),nil)
	} else {
		c.SetSession("uid", account.Id)
		models.DefaultAccountLoginLog.AccountLogin(&account)
		token ,tokenerr:=models.DefaultAccountToken.GenAccounttoken(&account)
		if(tokenerr!=nil){
			c.ErrorJson(20211201164342,tokenerr.Error(),nil)
		}
		accountRes :=AccountResponse{Id:account.Id,
			Name:account.Name,
			Email: account.Email,
			Token:token,
		}

		c.SuccessJson(accountRes)
	}
}


