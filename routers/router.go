package routers

import (
	"marketing/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/campaign/create", &controllers.CampaignController{},"post:CreateCampaign")
}
