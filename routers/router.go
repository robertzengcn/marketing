package routers

import (
	"fmt"
	"strings"
	"marketing/controllers"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

func init() {
	langs, err := beego.AppConfig.String("langs")  // 1
	if err != nil {  // 2
  	fmt.Println("Failed to load languages from the app.conf")
  	// return
	}
	langsArr := strings.Split(langs, "|")
	for _, lang := range langsArr {  // 4
		if err := i18n.SetMessage(lang, "conf/"+lang+".ini"); err != nil {  // 5
		  fmt.Println("Failed to set message file for l10n")
		//   return
		}
	  }
	  
    beego.Router("/", &controllers.MainController{})
	beego.Router("/campaign/create", &controllers.CampaignController{},"post:CreateCampaign")
	beego.Router("/login/accountlogin", &controllers.AccountController{},"post:Validaccount")
	beego.Router("/campaign", &controllers.CampaignController{},"get:ListCampaign")
}
