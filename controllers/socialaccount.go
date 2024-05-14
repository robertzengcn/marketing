package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"marketing/dto"
	"marketing/models"
	"marketing/utils"
	"strings"
	"strconv"
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
	social_type_id, _ := c.GetInt64("social_type_id")
	id, _ := c.GetInt64("id", 0) //social account id
	user := c.GetString("user")
	pass := c.GetString("pass")
	name := c.GetString("name")
	phone := c.GetString("phone")
	email := c.GetString("email")

	// proxyId, _ := c.GetInt64("proxy_id")

	status, _ := c.GetInt8("status", 1)
	// if proxyId > 0 {
	// 	//check proxy id valid
	// 	proxyModel := models.Proxy{}
	// 	proxy, _ := proxyModel.GetProxyById(proxyId)
	// 	if proxy.Id <= 0 {
	// 		c.ErrorJson(20240305105533, "proxy id incorrect", nil)
	// 	}
	// }

	if social_type_id <= 0 {
		c.ErrorJson(202304041002135, "social type incorrect", nil)
	}
	if len(user) <= 0 {
		c.ErrorJson(202304041005138, "user incorrect", nil)
	}
	if len(pass) <= 0 {
		c.ErrorJson(202304041005141, "pass incorrect", nil)
	}
	if len(email) > 0 {
		//check email vaild
		if !utils.ValidEmail(email) {
			c.ErrorJson(202304041115158, "email incorrect", nil)
		}
	}
	//config proxy
	var proxys []int64
	//update account proxy
	inputValues, _ := c.Input()
	for k, v := range inputValues {
		if k == "proxy[]" {
			if len(v) > 0 {
				// Convert v from []string to []int64
				proxyIDs := make([]int64, len(v))
				for i, val := range v {
					proxyID, _ := strconv.ParseInt(val, 10, 64)
					proxyIDs[i] = proxyID
				}
				proxys = append(proxys, proxyIDs...)
			}
		}
	}

	//check social type id correct
	socialplatform, err := models.DefaultSocialPlatform.GetSocialPlatformById(social_type_id)
	if err != nil {
		c.ErrorJson(202304041013163, err.Error(), nil)
	}
	var finId int64 = 0

	//check social account exist
	soa := models.SocialAccountUpdate{
		AccountId:      accountId,
		UserName:       user,
		PassWord:       pass,
		Socialplatform: &models.SocialPlatform{Id: socialplatform.Id},
		Status:         status,
		// Proxy:          &models.SocialProxy{Id: proxyId},
		AccountName:    name,
		Phone:          phone,
		Email:          email,
	}
	if id > 0 { //update exist social account
		finId = id
		err := models.DefaultSocialAccount.UpdateSocialAccount(id, accountId, soa)
		if err != nil {
			c.ErrorJson(20240312103083, err.Error(), nil)
		}
	} else {
		//save NEW social account
		saId, err := models.DefaultSocialAccount.CreateSocialAccount(soa)
		if err != nil {
			c.ErrorJson(20230403105183, err.Error(), nil)
		}
		finId = saId
	}
	socialAccountlistModel:=models.SocialAccountProxyList{}
	socialAccountlistModel.UpdateProxysToSocialAccount(finId,proxys)
	sap := dto.SocialAccountSaveresp{Id: finId}
	c.SuccessJson(sap)
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
	//get social proxy list
	// sopmodel := models.SocialProxy{}
	// logs.Info(socialaccounts)
	// socialproxy, err := sopmodel.GetSocialProxyById(socialaccounts.Proxy.Id)
	soaproxylist := models.SocialAccountProxyList{}
	socialproxy, err := soaproxylist.GetProxyBySocialAccountId(socialaccounts.Id)

	if err != nil && (err.Error() != "<QuerySeter> no row found") {
		c.ErrorJson(20230403094479, err.Error(), nil)
	}
	
	var sopArr []dto.SocialProxyDto
		for _, v := range socialproxy {
		sop := dto.SocialProxyDto{
			Id:       v.Id,
			Protocol: v.Protocol,
			Host:    v.Host,
			Port:    v.Port,
			// Url:      v.Protocol+"://"+v.Host + ":" + v.Port,
			Username: v.User,
			Password: v.Pass,
		}
		sopArr = append(sopArr, sop)
	}


	socialplatform := models.SocialPlatform{}
	socialplatform, err = socialplatform.GetSocialPlatformById(socialaccounts.Socialplatform.Id)
	if err != nil {
		c.ErrorJson(202304051034123, err.Error(), nil)
	}
	//get social profile
	socialprofile := models.SocialProfile{}
	profile, perr := socialprofile.GetSocialProfileByAccountId(socialaccounts.Id)
	if perr != nil {
		logs.Error(perr)
	}
	logs.Info(profile)
	socirep := dto.SocialAccountDetail{
		Id:            socialaccounts.Id,
		SocialType:    socialplatform.Name,
		SocialTypeId:  socialplatform.Id,
		SocialTypeUrl: socialplatform.Url,
		User:          socialaccounts.UserName,
		Pass:          socialaccounts.PassWord,
		Status:        socialaccounts.Status,
		Name:          profile.Name,
		PhoneNumber:   profile.PhoneNumber,
		Email:         profile.Email,
		Proxy:         sopArr,
	}
	c.SuccessJson(socirep)
}

//list social account
func (c *SocialAccountController) Listsocialaccount() {
	//get social account list from db
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	keyword := c.GetString("search")
	start, _ := c.GetInt("page", 0)
	size, _ := c.GetInt("size", 10)
	orderby := c.GetString("orderby")
	neworderby := strings.ReplaceAll(orderby, "-", "")
	if len(neworderby) > 0 {

		orderbyvaild := utils.Contains([]string{"id", "user_name"}, neworderby)
		if !orderbyvaild {
			c.ErrorJson(202403101644189, "orderby incorrect", nil)
		}
	}

	platformId, _ := c.GetInt64("platform", 0)
	san, sanerr := models.DefaultSocialAccount.CountSocialaccount(accountId, keyword, platformId)
	if sanerr != nil {
		c.ErrorJson(202403080943186, sanerr.Error(), nil)
	}

	socialaccounts, err := models.DefaultSocialAccount.ListSocialaccount(accountId, start, size, keyword, platformId, orderby)
	if err != nil {
		c.ErrorJson(202403060939156, err.Error(), nil)
	}
	var socialaccDto []dto.SocialAccountDto
	soaproxylist := models.SocialAccountProxyList{}
	for _, v := range socialaccounts {
		//get proxy by social account id
		// sopmodel := models.SocialProxy{}
		var useProxy int8 = 0
		socialproxy, _ := soaproxylist.GetProxyBySocialAccountId(v.Id)
		if len(socialproxy) > 0 {
			useProxy = 1
		}
		socialaccDto = append(socialaccDto, dto.SocialAccountDto{
			Id:           v.Id,
			SocialType:   v.Socialplatform.Name,
			SocialTypeId: v.Socialplatform.Id,
			User:         v.UserName,
			Password:     v.PassWord,
			Status:       v.Status,
			UseProxy:    useProxy,
		})
		// logs.Info(i,v)
	}
	salresp := dto.SocialAccountListresp{
		Total:   san,
		Records: socialaccDto,
	}
	c.SuccessJson(salresp)
}

//delete social account
func (c *SocialAccountController) DeleleSocialaccount() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	id, _ := c.GetInt64("id", 0)
	if id <= 0 {
		c.ErrorJson(202403101644189, "id incorrect", nil)
	}
	err := models.DefaultSocialAccount.DeleteSocialAccount(id, accountId)
	if err != nil {
		c.ErrorJson(202403080943186, err.Error(), nil)
	}
	c.SuccessJson(nil)
}
