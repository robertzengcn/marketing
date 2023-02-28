package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
)

type SocialController struct {
	BaseController
	i18n.Locale
}

func (c *SocialController) ChildPrepare(){
	// l := logs.GetLogger()
    //      l.Println("22222")
}

func (c *SocialController) CreateSocialAccount() {
	// l := logs.GetLogger()
	// l.Println("this is a message of get create social account")
	
	social_name := c.GetString("social_name")
	social_phone := c.GetString("social_phone")
	social_platformid,serr := c.GetInt64("social_type")
	if(serr!=nil){
		c.ErrorJson(202302211119,serr.Error(),nil)
	}
	//social_url := c.GetString("social_url")
	social_username := c.GetString("social_username")
	social_password := c.GetString("social_password")
	//social_token := c.GetString("social_token")
	social_email := c.GetString("social_email")
	campaign_id,_ := c.GetInt64("campaign_id")

	social_id,err:=models.DefaultSocialAccount.CreateSocialAccount(campaign_id,social_username,social_password,social_platformid,social_name,social_phone,social_email)

	if err != nil {
		c.ErrorJson(20211117161926,err.Error(),nil)
	}

	c.SuccessJson(social_id)
}




