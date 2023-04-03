package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	// "marketing/utils"
	"github.com/beego/i18n"
	"marketing/utils"
	// "fmt"
	//"errors"
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
	cts:=c.GetString("type","email")
	typearr:=[]string{"email","social"}
	if(!utils.Contains(typearr,cts)){
		c.ErrorJson(20230403100448,"type incorrect",nil)	
	}
	campaing_id,err:=models.DefaultCampaign.CreateCampaign(campaign_name,tag,cts)

	if err != nil {
		c.ErrorJson(20211117161926,err.Error(),nil)
	}
	if(cts=="social"){
		//for social type, we need to save proxy
		proxyUrl := c.GetString("proxy_url")
		proxyUser:=c.GetString("proxy_user")
		proxyPass:=c.GetString("proxy_pass")
		sop:=models.SocialProxy{
		Url:proxyUrl,
		Username: proxyUser,
		Password: proxyPass,
		Campaign: &models.Campaign{CampaignId: campaing_id},

		}

		socialProxyM:=models.SocialProxy{}
		//valid proxy data before save data
		
		_,spr:=socialProxyM.Save(sop)
		if(spr!=nil){
			c.ErrorJson(20230403105774,spr.Error(),nil)
		}
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
	//get social proxy
	sopmodel:=models.SocialProxy{}
	socialproxy,err:=sopmodel.GetSocialProxyByCampaignId(campaign_id)
	if(err!=nil){
		c.ErrorJson(20230403094479,err.Error(),nil)
	}
	sop:=SoProxy{
		Url: socialproxy.Url,
		User: socialproxy.Username,
		Pass: socialproxy.Password,
	}
	socirep:=Socialresp{
		User: socialaccounts.UserName,
		Pass: socialaccounts.PassWord,
		Sotype: socialaccounts.SocialplatformId.Name,
		Proxy: sop,
	}
	c.SuccessJson(socirep)
}



