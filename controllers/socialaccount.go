package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"marketing/dto"
	"marketing/models"
	"marketing/utils"
)

type SocialAccountController struct {
	BaseController
	i18n.Locale
}

func (c *SocialAccountController) ChildPrepare() {
	// l := logs.GetLogger()
	//      l.Println("22222")
}

//save social account
func (c *SocialAccountController) Savesocialaccount() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	social_type := c.GetString("social_type")
	user := c.GetString("user")
	pass := c.GetString("pass")
	accountname := c.GetString("accountname")
	phone := c.GetString("phone")
	email := c.GetString("email")
	proxyId, _ := c.GetInt64("proxy_id")
	status, _ := c.GetInt8("status", 1)
	if proxyId > 0 {
		//check proxy id valid
		proxyModel := models.Proxy{}
		proxy, _ := proxyModel.GetProxyById(proxyId)
		if proxy.Id <= 0 {
			c.ErrorJson(20240305105533, "proxy id incorrect", nil)
		}
	}
	// proxyUrl := c.GetString("proxy_url", "")
	// proxyUser := c.GetString("proxy_user", "")
	// proxyPass := c.GetString("proxy_pass", "")

	if len(social_type) <= 0 {
		c.ErrorJson(202304041002135, "social type incorrect", nil)
	}
	if len(user) <= 0 {
		c.ErrorJson(202304041005138, "user incorrect", nil)
	}
	if len(pass) <= 0 {
		c.ErrorJson(202304041005141, "pass incorrect", nil)
	}
	// if(len(proxyUrl)<=0){
	// 	c.ErrorJson(202304041010144,"proxy url incorrect",nil)
	// }
	// if(len(proxyUser)<=0){
	// 	c.ErrorJson(202304041012147,"proxy user incorrect",nil)
	// }
	// if(len(proxyPass)<=0){
	// 	c.ErrorJson(202304041012150,"proxy pass incorrect",nil)
	// }
	//check email vaild
	if !utils.ValidEmail(email) {
		c.ErrorJson(202304041115158, "email incorrect", nil)
	}

	//check campaign id correct
	// campaign, err := models.DefaultCampaign.FindCambyid(campaign_id)
	// if err != nil {
	// 	c.ErrorJson(202304041037159, err.Error(), nil)
	// }

	// var sop models.SocialProxy
	// if len(proxyUrl) > 0 {
	// 	sop = models.SocialProxy{
	// 		Url:      proxyUrl,
	// 		Username: proxyUser,
	// 		Password: proxyPass,
	// 	}
	// }
	//check social platform id
	socialplatform, err := models.DefaultSocialPlatform.FindSocialPlatformByName(social_type)
	if err != nil {
		c.ErrorJson(202304041013163, err.Error(), nil)
	}
	// var sproxyId int64
	// var spr error
	// if sop != (models.SocialProxy{}) {
	// 	socialProxyM := models.SocialProxy{}
	// 	//valid proxy data before save data
	// 	sproxyId, spr = socialProxyM.Save(sop)
	// 	if spr != nil {
	// 		c.ErrorJson(20230403105774, spr.Error(), nil)
	// 	}
	// }
	soa:=models.SocialAccountUpdate{
		AccountId: accountId,
	UserName:user,        
	PassWord:pass,     
	Socialplatform: &models.SocialPlatform{Id:socialplatform.Id},
	Status:status,                     
	Proxy:           &models.SocialProxy{Id:proxyId},
	AccountName:accountname,  
	Phone:		   phone,
	Email:email,
	}
	//save social account
	socialaccount, err := models.DefaultSocialAccount.CreateSocialAccount(soa)
	if err != nil {
		c.ErrorJson(20230403105183, err.Error(), nil)
	}
	c.SuccessJson(socialaccount)
}

//get social account by id
func (c *SocialAccountController) GetSocialAccount() {
	id, _ := c.GetInt64("id", 0)
	if id <= 0 {
		c.ErrorJson(202304140951101, "id incorrect", nil)
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	// logs.Info("campaign id",campaign_id)
	socialaccounts, err := models.DefaultSocialAccount.GetSocialAccount(id, accountId)
	if err != nil {

		c.ErrorJson(202304050957100, err.Error(), nil)
	}
	//get social proxy
	sopmodel := models.SocialProxy{}
	logs.Info(socialaccounts)
	socialproxy, err := sopmodel.GetSocialProxyById(socialaccounts.Proxy.Id)

	if err != nil && (err.Error() != "<QuerySeter> no row found") {
		c.ErrorJson(20230403094479, err.Error(), nil)
	}
	// var sop SoProxy
	// if socialproxy != (models.SocialProxy{}) {
		sop := dto.SocialProxyDto{
			Url:  socialproxy.Url,
			Username: socialproxy.Username,
			Password: socialproxy.Password,
		}

	// }
	// logs.Info(socialaccounts.Socialplatform.Id)
	//get social platform name
	socialplatform := models.SocialPlatform{}
	socialplatform, err = socialplatform.GetSocialPlatformById(socialaccounts.Socialplatform.Id)
	if err != nil {
		c.ErrorJson(202304051034123, err.Error(), nil)
	}
	//get social profile
	socialprofile:=models.SocialProfile{}
	profile,perr:=socialprofile.GetSocialProfileByAccountId(socialaccounts.Id)
	if(perr!=nil){
		logs.Error(perr)
	}
	logs.Info(profile)
	socirep := dto.SocialAccountDetail{
		Id:		   socialaccounts.Id,
		SocialType: socialplatform.Name,
		SocialTypeId: socialplatform.Id,
		User:   socialaccounts.UserName,
		Pass:   socialaccounts.PassWord,
		Status: socialaccounts.Status,
		Name:   profile.Name,
		PhoneNumber: profile.PhoneNumber,
		Email: profile.Email,
		Proxy:  sop,
	}
	c.SuccessJson(socirep)
}

//list social account
func (c *SocialAccountController) Listsocialaccount() {
	//get social account list from db
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	start, _ := c.GetInt("page", 0)
	// logs.Info("start is "+strconv.Itoa(start))
	size, _ := c.GetInt("size", 10)
	socialaccounts, err := models.DefaultSocialAccount.ListSocialaccount(accountId, start, size)
	if err != nil {
		c.ErrorJson(202403060939156, err.Error(), nil)
	}
	var socialaccDto []dto.SocialAccountDto
	for _, v := range socialaccounts {
		socialaccDto = append(socialaccDto, dto.SocialAccountDto{
			Id:           v.Id,
			SocialType:   v.Socialplatform.Name,
			SocialTypeId: v.Socialplatform.Id,
			User:         v.UserName,
			Password:     v.PassWord,
			Status:       v.Status,
			// Proxy: dto.ProxyDto{
			// 	Id:       v.Proxy.Id,
			// 	Url:      v.Proxy.Url,
			// 	Username: v.Proxy.Username,
			// 	Password: v.Proxy.Password,
			// },
		})
		// logs.Info(i,v)
	}
	c.SuccessJson(socialaccDto)
}

//update social account
func (c *SocialAccountController) Updatesocialaccount() {
	socialaccount_id, _ := c.GetInt64("socialaccount_id")
	if socialaccount_id <= 0 {
		c.ErrorJson(202403061055157, "social account id incorrect", nil)
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	user := c.GetString("user")
	pass := c.GetString("pass")
	accountname := c.GetString("accountname")
	phone := c.GetString("phone")
	email := c.GetString("email")
	proxyId, _ := c.GetInt64("proxy_id")
	status, _ := c.GetInt8("status", 0)
	if status != 0 && status != 1 {
		c.ErrorJson(202403061052195, "status incorrect", nil)
	}
	social_type_id, serr := c.GetInt64("social_type_id")
	if serr != nil {
		c.ErrorJson(202403061022195, "social type id incorrect", nil)
	}
	if social_type_id <= 0 {
		c.ErrorJson(202403061022198, "social type incorrect", nil)
	}
	SocialPlatform, err := models.DefaultSocialPlatform.GetSocialPlatformById(social_type_id)
	if SocialPlatform.Id <= 0 || err != nil {
		c.ErrorJson(202304041013163, err.Error(), nil)
	}
	if proxyId > 0 {
		//check proxy id valid
		proxyModel := models.Proxy{}
		proxy, _ := proxyModel.GetProxyById(proxyId)
		if proxy.Id <= 0 {
			c.ErrorJson(20240305105533, "proxy id incorrect", nil)
		}
	}
	//check email vaild
	if !utils.ValidEmail(email) {
		c.ErrorJson(202304041115158, "email incorrect", nil)
	}
	socialAccountmodel := models.SocialAccount{}
	socialUpdate := models.SocialAccountUpdate{
		UserName:user,
		PassWord:pass,
		Socialplatform: &models.SocialPlatform{Id: SocialPlatform.Id},
		Status:status,
		Proxy: & models.SocialProxy{Id: proxyId},
	}
	err = socialAccountmodel.UpdateSocialAccount(socialaccount_id, accountId, socialUpdate)
	if err != nil {
		c.ErrorJson(202403061057198, err.Error(), nil)
	}
	//update social profile
	socialProfilemodel := models.SocialProfile{}
	profileEntity:=models.SocialProfileUpdate{
		Name:accountname,
		PhoneNumber: phone,
		Email:email,
	}
	_,surr:=socialProfilemodel.UpdateSocialProfile(socialaccount_id, profileEntity)
	if(surr!=nil){
		c.ErrorJson(202403061501246, surr.Error(), nil)
	}
	c.SuccessJson("update success")
}
