package controllers

import (
	"marketing/models"
	// "errors"
	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/i18n"
	"marketing/dto"
	"strings"
	"marketing/utils"
)

type EmailtplController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}
type CreateEmailres struct {
	Id int64
}

/// create email template
func (c *EmailtplController) CreateEmailtpl() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	c.Lang = c.BaseController.Lang
	email_title := c.GetString("email_title")
	email_content := c.GetString("email_content")
	campaign_id, _ := c.GetInt64("campaign_id")
	// if cerr != nil {
	// 	c.ErrorJson(20240925103730, cerr.Error(), nil)
	// }
	campaiginVar := &models.Campaign{}
	if campaign_id > 0 {
		CampaignModel := models.Campaign{}

		campaiginVar, camerr := CampaignModel.FindCambyid(campaign_id)
		if camerr != nil {
			c.ErrorJson(20220618162729, camerr.Error(), nil)
		}
		if campaiginVar == nil {

			c.ErrorJson(20220618162832, c.Tr("invail_campaign_id"), nil)
		}
	}
	// logs.Info(campaiginVar)
	emailtplModel := models.EmailTpl{}
	emailVar := models.EmailTpl{
		TplTitle:   email_title,
		TplContent: email_content,
		CampaignId: campaiginVar,
		Status:     1,
		AccountId:  &models.Account{Id: accountId},
	}

	valid := validation.Validation{}
	b, verr := valid.Valid(&emailVar)
	if !b {
		var errMessage string
		for _, err := range valid.Errors {
			errMessage += err.Key + ":" + err.Message
		}
		c.ErrorJson(202409241405104, errMessage, nil)
	}
	if verr != nil {
		c.ErrorJson(202409241405104, verr.Error(), nil)
	}
	// logs.Info(emailVar)
	emailId, emailerr := emailtplModel.Createone(emailVar)
	if emailerr != nil {
		c.ErrorJson(20220617155924, emailerr.Error(), nil)
	}
	c.SuccessJson(CreateEmailres{Id: emailId})
}

///get email list by account id
func (c *EmailtplController) GetEmailtplList() {
	page, perr := c.GetInt64("page", 0)
	if perr != nil {
		c.ErrorJson(20240924112459, perr.Error(), nil)
	}
	size, serr := c.GetInt64("size", 10)
	if serr != nil {
		c.ErrorJson(20240924112463, serr.Error(), nil)
	}
	search:=c.GetString("search","")

	orderby := c.GetString("orderby","")
	neworderby := strings.ReplaceAll(orderby, "-", "")
	if len(neworderby) > 0 {

		orderbyvaild := utils.Contains([]string{"tpl_id", "tpl_record"}, neworderby)
		if !orderbyvaild {
			c.ErrorJson(202403101644189, "orderby incorrect", nil)
		}
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailtplModel := models.EmailTpl{}
	emailtplList, emailerr := emailtplModel.GetEmailTplListByAccountId(accountId, page, size,search,neworderby)
	if emailerr != nil {
		c.ErrorJson(20220617155924, emailerr.Error(), nil)
	}
	var emailtpldtos []dto.EmailtplDto
	for i, s := range emailtplList {
		emailtpldtos = append(emailtpldtos, dto.EmailtplDto{
			TplId:      s.TplId,
			TplTitle:   s.TplTitle,
			TplContent: s.TplContent,
			// CampaignId: s.CampaignId.CampaignId,
			TplRecord:  s.TplRecord.Format("2006-01-02 15:04:05"),
			Status:     s.Status,
			Index:int64(page)+int64(i)+1,
		})
	}
	ecnum,ecerr:=emailtplModel.GetEmailTplCountByAccountId(accountId,search)
	if ecerr != nil {
		c.ErrorJson(202409271129118, ecerr.Error(), nil)
	}

	// Removed unused variable resp and fixed instantiation issue
	resp := dto.CommonResponse[[]dto.EmailtplDto]{
		Record: emailtpldtos,
		Total:    ecnum,
	}

	c.SuccessJson(resp)
	// c.SuccessJson(emailtpldtos)
}

/// get email template by account id and tpl id
func (c *EmailtplController) GetEmailtplById() {
	tplId, terr := c.GetInt64(":id")
	if terr != nil {
		c.ErrorJson(20220618162124, terr.Error(), nil)
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailtplModel := models.EmailTpl{}
	emailtpl, emailerr := emailtplModel.GetEmailTplByIdAndAccountId(tplId, accountId)
	if emailerr != nil {
		c.ErrorJson(20220617155924, emailerr.Error(), nil)
	}
	// var emailtpldto dto.EmailtplDto
	emailtpldto := dto.EmailtplDto{
		TplId:      emailtpl.TplId,
		TplTitle:   emailtpl.TplTitle,
		TplContent: emailtpl.TplContent,
		// CampaignId: emailtpl.CampaignId.CampaignId,
		TplRecord:  emailtpl.TplRecord.Format("2006-01-02 15:04:05"),
		Status:     emailtpl.Status,
	}
	c.SuccessJson(emailtpldto)
}

///update email template by id
func (c *EmailtplController) UpdateEmailtpl() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	tplId, terr := c.GetInt64(":id")
	if terr != nil {
		c.ErrorJson(20220618162124, terr.Error(), nil)
	}
	email_title := c.GetString("email_title")
	email_content := c.GetString("email_content")
	campaign_id, cerr := c.GetInt64("campaign_id")
	status, sterr := c.GetInt("status", 1)
	if sterr != nil {
		c.ErrorJson(202409241405104, sterr.Error(), nil)
	}
	if cerr != nil {
		c.ErrorJson(20220618162124, cerr.Error(), nil)
	}
	campaiginVar := &models.Campaign{}
	if campaign_id > 0 {
		CampaignModel := models.Campaign{}
		campaiginVar, camerr := CampaignModel.FindCambyid(campaign_id)
		if camerr != nil {
			c.ErrorJson(20220618162729, camerr.Error(), nil)
		}
		if campaiginVar == nil {

			c.ErrorJson(20220618162832, c.Tr("invail_campaign_id"), nil)
		}

	}

	emailtplModel := models.EmailTpl{}
	eres, ereserr := emailtplModel.GetEmailTplByIdAndAccountId(tplId, accountId)
	if ereserr != nil {
		c.ErrorJson(202409241410128, ereserr.Error(), nil)
	}
	if eres == nil {
		c.ErrorJson(202409241410131, "email template not exist", nil)
	}

	emailVar := models.EmailTpl{
		TplId:      tplId,
		TplTitle:   email_title,
		TplContent: email_content,
		CampaignId: campaiginVar,
		Status:     status,
		AccountId: &models.Account{Id: accountId},
	}
	valid := validation.Validation{}
	b, verr := valid.Valid(&emailVar)
	if !b {
		var errMessage string
		for _, err := range valid.Errors {
			errMessage += err.Key + ":" + err.Message
		}
		c.ErrorJson(202409241405104, errMessage, nil)
	}
	if verr != nil {
		c.ErrorJson(202409241405104, verr.Error(), nil)
	}
	_, err := emailtplModel.UpdateEmailTplById(emailVar)
	if err != nil {
		c.ErrorJson(202409241405104, err.Error(), nil)
	}
	c.SuccessJson(CreateEmailres{Id: tplId})

}
///delete email template by id and account id
func (c *EmailtplController) DeleteEmailtpl() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	tplId, terr := c.GetInt64(":id")
	if terr != nil {
		c.ErrorJson(20220618162124, terr.Error(), nil)
	}
	emailtplModel := models.EmailTpl{}
	eres, ereserr := emailtplModel.GetEmailTplByIdAndAccountId(tplId, accountId)
	if ereserr != nil {
		c.ErrorJson(202409241410128, ereserr.Error(), nil)
	}
	if eres == nil {
		c.ErrorJson(202409241410131, "email template not exist", nil)
	}
	res, err := emailtplModel.DeleteEmailTplByIdAndAccountId(tplId,accountId)
	if err != nil {
		c.ErrorJson(202409241405104, err.Error(), nil)
	}
	if(res<=0){
		c.ErrorJson(202409301448246, "delete fail", nil)
	}
	c.SuccessJson(tplId)
}
