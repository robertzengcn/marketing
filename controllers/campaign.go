package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	// "marketing/utils"
	"github.com/beego/i18n"
	// "fmt"
)

type CampaignController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

type Socialresp struct{
	Sotype string `json:"sotype"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Proxy SoProxy `json:"proxy"`
}
type SoProxy struct{
	Url string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}
// func (c *CampaignController) Prepare() {
//     c.EnableXSRF = false
// }
func (c *CampaignController)ChildPrepare(){
	
}


//create campaign
func (c *CampaignController) CreateCampaign() {
	// l := logs.GetLogger()
	// l.Println("this is a message of get create campaign")
	
	campaign_name := c.GetString("campaign_name")
	tag := c.GetString("tag")

	campaing_id,err:=models.DefaultCampaign.CreateCampaign(campaign_name,tag)

	if err != nil {
		c.ErrorJson(20211117161926,err.Error(),nil)
	}

	c.SuccessJson(campaing_id)
}
//list campaign use request
func (c *CampaignController) ListCampaign() {
	start,_ := c.GetInt("start",0)
	num,_:= c.GetInt("number",10)
	
	campagins,err:=models.DefaultCampaign.ListCampaign(start,num)
	if(err!=nil){
		c.ErrorJson(20211208153839,err.Error(),nil)

	}
	c.SuccessJson(campagins)
}

//get socail account relation with campaign use campaign Id
func (c *CampaignController) GetSocialAccount() {
	campaign_id,_ := c.GetInt64("campaign_id",0)
	
	socialaccounts,err:=models.DefaultCampaign.GetSocialAccount(campaign_id)
	if(err!=nil){
		c.ErrorJson(20211208153839,err.Error(),nil)

	}
	socirep:=Socialresp{
		User: socialaccounts.UserName,
		Pass: socialaccounts.PassWord,
		Sotype: socialaccounts.SocialplatformId.Name,
		
	}
	c.SuccessJson(socirep)
}



