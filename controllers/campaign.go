package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	// "marketing/utils"
	// "marketing/utils"
	"marketing/dto"
	// "github.com/beego/beego/v2/core/logs"
	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	// "fmt"
	//"errors"
	// "strconv"
)

type CampaignController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

type Socialresp struct {
	Sotype string  `json:"sotype"`
	User   string  `json:"user"`
	Pass   string  `json:"pass"`
	Proxy  SoProxy `json:"proxy"`
}
type SoProxy struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type CampaignlistResponse struct {
	Data []models.Campaign `json:"data"`
	Num  int64             `json:"num"`
}

// func (c *CampaignController) Prepare() {
//     c.EnableXSRF = false
// }
func (c *CampaignController) ChildPrepare() {

}

//create campaign
func (c *CampaignController) CreateCampaign() {
	// l := logs.GetLogger()
	// l.Println("this is a message of get create campaign")

	campaign_name := c.GetString("campaign_name")
	tag := c.GetString("tag")
	// cts := c.GetString("type", "email")
	// typearr := []string{"email", "social"}
	// if !utils.Contains(typearr, cts) {
	// 	c.ErrorJson(20230403100448, "type incorrect", nil)
	// }
	cts,cterr:=c.GetInt16("type")
	if(cterr!=nil){
		c.ErrorJson(20231128093656, cterr.Error(), nil)
	}
	if(cts<=0){
		c.ErrorJson(20231128093759, "type incorrect", nil)
	}

	campaing_id, err := models.DefaultCampaign.CreateCampaign(campaign_name, tag, cts)

	if err != nil {
		c.ErrorJson(20211117161926, err.Error(), nil)
	}
	// if(cts=="social"){

	// 	//for social type, we need to save proxy
	// 	proxyUrl := c.GetString("proxy_url")
	// 	proxyUser:=c.GetString("proxy_user")
	// 	proxyPass:=c.GetString("proxy_pass")
	// 	sop:=models.SocialProxy{
	// 	Url:proxyUrl,
	// 	Username: proxyUser,
	// 	Password: proxyPass,
	// 	//Campaign: &models.Campaign{CampaignId: campaing_id},

	// 	}

	// 	socialProxyM:=models.SocialProxy{}
	// 	//valid proxy data before save data

	// 	_,spr:=socialProxyM.Save(sop)
	// 	if(spr!=nil){
	// 		c.ErrorJson(20230403105774,spr.Error(),nil)
	// 	}
	// }

	c.SuccessJson(campaing_id)
}

//list campaign use request
func (c *CampaignController) ListCampaign() {
	
	start, _ := c.GetInt("page", 0)
	// logs.Info("start is "+strconv.Itoa(start))
	num, _ := c.GetInt("size", 10)
	uid := c.GetSession("uid")
	accountId:=uid.(int64)
	campagins, err := models.DefaultCampaign.ListCampaign(start, num,accountId)
	if err != nil {
		c.ErrorJson(20211208153839, err.Error(), nil)
	}
	//count compaign number
	campaginNum,camNumerr:=models.DefaultCampaign.CountCampaign()
	if(camNumerr!=nil){
		c.ErrorJson(20211208153942, camNumerr.Error(), nil)
	}
	camaigndto:=dto.CampaignDto{
		Num: campaginNum,
	}
	for _, element := range campagins {
		// logs.Info(element)
		camaigndto.Records=append(camaigndto.Records,dto.CampaignItemDto{
			CampaignId:element.CampaignId,
			CampaignName:element.CampaignName,
			CampaignDescription:element.CampaignDescription,
			Disable:element.Disable,
			Types:element.Types.CampaignTypeName,
			// Tags:element.Tags,
			AccountId: element.AccountId.Id,
		})
	}
	// CampaignlistResponse := CampaignlistResponse{
	// 	Data: campagins,
	// 	Num:  campaginNum,
	// }
	c.SuccessJson(camaigndto)
}

//get socail account relation with campaign use campaign Id
// func (c *CampaignController) GetSocialAccount() {
// 	campaign_id, _ := c.GetInt64("campaign_id", 0)
// 	if campaign_id <= 0 {
// 		c.ErrorJson(202304140951101, "campaign id incorrect", nil)
// 	}
// 	// logs.Info("campaign id",campaign_id)
// 	socialaccounts, err := models.DefaultSocialAccount.GetSocialAccountcam(campaign_id)
// 	if err != nil {
		
// 		c.ErrorJson(202304050957100, err.Error(), nil)
// 	}
// 	//get social proxy
// 	sopmodel := models.SocialProxy{}
// 	logs.Info(socialaccounts)
// 	socialproxy, err := sopmodel.GetSocialProxyById(socialaccounts.Proxy.Id)

// 	if err != nil &&(err.Error() !="<QuerySeter> no row found"){
// 		c.ErrorJson(20230403094479, err.Error(), nil)
// 	}
// 	var sop SoProxy
// 	if(socialproxy!=(models.SocialProxy{})){
// 	sop = SoProxy{
// 		Url:  socialproxy.Url,
// 		User: socialproxy.Username,
// 		Pass: socialproxy.Password,
// 	}
// 	}
// 	logs.Info(socialaccounts.SocialplatformId.Id)
// 	//get social platform name
// 	socialplatform := models.SocialPlatform{}
// 	socialplatform, err = socialplatform.GetSocialPlatformById(socialaccounts.SocialplatformId.Id)
// 	if err != nil {
// 		c.ErrorJson(202304051034123, err.Error(), nil)
// 	}
// 	logs.Info(socialplatform)
// 	socirep := Socialresp{
// 		User:   socialaccounts.UserName,
// 		Pass:   socialaccounts.PassWord,
// 		Sotype: socialplatform.Name,
// 		Proxy:  sop,
// 	}
// 	c.SuccessJson(socirep)
// }

//create social account in campaign
// func (c *CampaignController) CreateSocialAccount() {
// 	campaign_id, _ := c.GetInt64("campaign_id", 0)
// 	social_type := c.GetString("social_type")
// 	user := c.GetString("user")
// 	pass := c.GetString("pass")
// 	accountname := c.GetString("accountname")
// 	phone := c.GetString("phone")
// 	email := c.GetString("email")

// 	proxyUrl := c.GetString("proxy_url", "")
// 	proxyUser := c.GetString("proxy_user", "")
// 	proxyPass := c.GetString("proxy_pass", "")
// 	if campaign_id <= 0 {
// 		c.ErrorJson(202304041002132, "campaign id incorrect", nil)
// 	}
// 	if len(social_type) <= 0 {
// 		c.ErrorJson(202304041002135, "social type incorrect", nil)
// 	}
// 	if len(user) <= 0 {
// 		c.ErrorJson(202304041005138, "user incorrect", nil)
// 	}
// 	if len(pass) <= 0 {
// 		c.ErrorJson(202304041005141, "pass incorrect", nil)
// 	}
// 	// if(len(proxyUrl)<=0){
// 	// 	c.ErrorJson(202304041010144,"proxy url incorrect",nil)
// 	// }
// 	// if(len(proxyUser)<=0){
// 	// 	c.ErrorJson(202304041012147,"proxy user incorrect",nil)
// 	// }
// 	// if(len(proxyPass)<=0){
// 	// 	c.ErrorJson(202304041012150,"proxy pass incorrect",nil)
// 	// }
// 	//check email vaild
// 	if !utils.ValidEmail(email) {
// 		c.ErrorJson(202304041115158, "email incorrect", nil)
// 	}

// 	//check campaign id correct
// 	campaign, err := models.DefaultCampaign.FindCambyid(campaign_id)
// 	if err != nil {
// 		c.ErrorJson(202304041037159, err.Error(), nil)
// 	}

// 	var sop models.SocialProxy
// 	if len(proxyUrl) > 0 {
// 		sop = models.SocialProxy{
// 			Url:      proxyUrl,
// 			Username: proxyUser,
// 			Password: proxyPass,
// 		}
// 	}
// 	//check social platform id
// 	socialplatform, err := models.DefaultSocialPlatform.FindSocialPlatformByName(social_type)
// 	if err != nil {
// 		c.ErrorJson(202304041013163, err.Error(), nil)
// 	}
// 	var sproxyId int64
// 	var spr error
// 	if sop != (models.SocialProxy{}) {
// 		socialProxyM := models.SocialProxy{}
// 		//valid proxy data before save data
// 		sproxyId, spr = socialProxyM.Save(sop)
// 		if spr != nil {
// 			c.ErrorJson(20230403105774, spr.Error(), nil)
// 		}
// 	}
// 	//save social account
// 	socialaccount, err := models.DefaultSocialAccount.CreateSocialAccount(campaign.CampaignId, user, pass, socialplatform.Id, accountname, phone, email, sproxyId)
// 	if err != nil {
// 		c.ErrorJson(20230403105183, err.Error(), nil)
// 	}
// 	c.SuccessJson(socialaccount)
// }
