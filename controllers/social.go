package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	"strings"
)

type SocialController struct {
	BaseController
	i18n.Locale
}

func (c *SocialController) ChildPrepare() {
	// l := logs.GetLogger()
	//      l.Println("22222")
}

// func (c *SocialController) CreateSocialAccount() {
// 	// l := logs.GetLogger()
// 	// l.Println("this is a message of get create social account")

// 	social_name := c.GetString("social_name")
// 	social_phone := c.GetString("social_phone")
// 	social_platformid,serr := c.GetInt64("social_type")
// 	if(serr!=nil){
// 		c.ErrorJson(202302211119,serr.Error(),nil)
// 	}
// 	//social_url := c.GetString("social_url")
// 	social_username := c.GetString("social_username")
// 	social_password := c.GetString("social_password")
// 	//social_token := c.GetString("social_token")
// 	social_email := c.GetString("social_email")
// 	campaign_id,_ := c.GetInt64("campaign_id")

// 	social_id,err:=models.DefaultSocialAccount.CreateSocialAccount(campaign_id,social_username,social_password,social_platformid,social_name,social_phone,social_email,)

// 	if err != nil {
// 		c.ErrorJson(20211117161926,err.Error(),nil)
// 	}

// 	c.SuccessJson(social_id)
// }
type SocialEntity struct {
	CampaignId   int64    `json:"campaign_id"`
	CampaignName string   `json:"campaign_name"`
	Tags         string   `json:"tags"`
	Types        string   `json:"types"`
	Keyword      []string `json:"keyword"`
}
type SocialTaskEntity struct {
	Id           int64  `json:"id"`
	CampaignId   int64  `json:"campaign_id"`
	CampaignName string `json:"campaign_name"`
	Tags         string `json:"tag"`
	Type         string `json:"type"`
	Keywords     []string `json:"keywords"`
}

///list social campaign
func (c *SocialController) Listsocialcampaigin() {
	types := c.GetString("types", "social")
	start, sterr := c.GetInt("start", 0)
	if sterr != nil {
		c.ErrorJson(20230525093049, sterr.Error(), nil)
	}
	end, enderr := c.GetInt("end", 0)
	if enderr != nil {
		c.ErrorJson(20230525093053, sterr.Error(), nil)
	}
	campaignModel := models.Campaign{}
	campaigns, err := campaignModel.ListCambytype(types, start, end)
	if err != nil {
		c.ErrorJson(20230525093057, err.Error(), nil)
	}
	keywordModel := models.Keyword{}
	var socials []SocialEntity
	for i := 0; i < len(campaigns); i++ {
		var soentity SocialEntity
		soentity.CampaignId = campaigns[i].CampaignId
		soentity.CampaignName = campaigns[i].CampaignName
		soentity.Tags = campaigns[i].Tags
		soentity.Types = campaigns[i].Types
		tags := strings.Split(campaigns[i].Tags, ",")
		keywordArr, kErr := keywordModel.Getkeywordbytag(tags, 5)
		if kErr != nil {
			continue
		}
		if len(keywordArr) <= 0 {
			continue
		}
		// soentity.Keyword=[]string{}
		for _, s := range keywordArr {
			soentity.Keyword = append(soentity.Keyword, s.Keyword)
		}
		socials = append(socials, soentity)
	}
	c.SuccessJson(socials)

}

///get social task list
func (c *SocialController) Getsocialtasklist() {
	socialTaskmodel := models.SocialTask{}
	soTask, soterr := socialTaskmodel.GetSocialTaskList()
	if soterr != nil {
		c.ErrorJson(20230526095299, soterr.Error(), nil)
	}
	soentityArr := []SocialTaskEntity{}
	keywordModel := models.Keyword{}
	for _, stv := range soTask {
		var soentity SocialTaskEntity
		soentity.Id = stv.Id
		soentity.CampaignId = stv.Campaign.CampaignId
		soentity.CampaignName = stv.Campaign.CampaignName
		soentity.Tags = stv.Campaign.Tags
		soentity.Type = stv.Type
		tags := strings.Split(stv.Campaign.Tags, ",")
		keywordArr, kErr := keywordModel.Getkeywordbytag(tags, 5)
		if kErr != nil {
			continue
		}
		for _, s := range keywordArr {
			soentity.Keywords = append(soentity.Keywords, s.Keyword)
		}
		soentityArr = append(soentityArr, soentity)
	}
	c.SuccessJson(soentityArr)
}
