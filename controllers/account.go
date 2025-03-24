package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	//"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"marketing/models"
	// "strconv"
)

type AccountTokenResponse struct {
	// Id int64
	// Name    string
	// Email   string
	Token string
}
type AccountResp struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
type AccountController struct {
	BaseController
	i18n.Locale
}

// func (a *AccountController) Prepare() {
//     a.EnableXSRF = false
// }
func (c *AccountController) ChildPrepare() {
	// l := logs.GetLogger()
	//      l.Println("22222")
}

///valid account
func (c *AccountController) Validaccount() {
	username := c.GetString("username")
	pass := c.GetString("password")
	// logs.Info(username)
	// logs.Info(pass)
	if len(username) == 0 || len(pass) == 0 {
		//用户名和邮箱为空
		c.ErrorJson(20211122163020, c.Tr("email_pass_empty"), nil)
	}
	// l := logs.GetLogger()

	// l.Println(username)
	// l.Println(pass)
	// l.Println(c.Tr("welcome"))
	// l.Println("44444")
	accountModel := models.Account{}
	account, err := accountModel.Validaccount(username, pass)

	if err != nil {

		c.ErrorJson(202108031641149, err.Error(), nil)
	} else {
		c.SetSession("uid", account.Id)
		//record login log
		models.DefaultAccountLoginLog.AccountLogin(&account)
		token, tokenerr := models.DefaultAccountToken.GenAccounttoken(&account)
		if tokenerr != nil {
			c.ErrorJson(20211201164342, tokenerr.Error(), nil)
		}
		accountRes := AccountTokenResponse{
			// Id:account.Id,
			// Name:account.Name,
			// Email: account.Email,
			Token: token,
		}

		c.SuccessJson(accountRes)
	}

}

///echo user info by session
func (c *AccountController) Accountinfo() {
	uid := c.GetSession("uid")
	uidint64 := uid.(int64)
	// logs.Info("uid will be"+strconv.FormatInt(uidint64,10))
	if uid == nil {
		c.ErrorJson(202302270948, c.Tr("user_not_login"), nil)
	}
	accountModel := models.Account{}
	//convert uid to int64
	// uidint64:=uid.(int64)
	acc, aerr := accountModel.GetAccountbyid(uidint64)
	var accountrolearr []string
	if acc.Roles != nil {
		for _, element := range acc.Roles {
			accountrolearr = append(accountrolearr, element.Name)
		}

	}
	arsp := AccountResp{}
	arsp.Email = acc.Email
	arsp.Name = acc.Name
	arsp.Roles = accountrolearr
	if aerr != nil {
		c.ErrorJson(202302270950, aerr.Error(), nil)
	}
	c.SuccessJson(arsp)
}

func (c *AccountController) Signout() {
	tokenId := c.GetSession("tokenId")
	tokenId64 := tokenId.(int64)
	if tokenId64 > 0 {
		accountTokenModel:=models.AccountToken{}
		_,err:=accountTokenModel.DeleteAccountToken(tokenId64)
		if(err!=nil){
			c.ErrorJson(202403041013, err.Error(), nil)
		}
	}
	c.SuccessJson(nil)
}
