package controllers
import (
	"github.com/beego/i18n"
	"marketing/models"
	"marketing/utils"
	"marketing/dto"
)

type EmailserviceController struct {
	BaseController
	i18n.Locale
}

func (c *EmailserviceController) ChildPrepare(){

}

///add email service
func (c *EmailserviceController) Addemailservice(){
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	service_name := c.GetString("service_name")
	if(len(service_name)<1){
		c.ErrorJson(20220815102320, "get service name error", nil)
	}
	service_from := c.GetString("service_from")
	if(len(service_from)<1){
		c.ErrorJson(20220815100319, "get email service from error", nil)
	}
	if(!utils.ValidEmail(service_from)){
		c.ErrorJson(20220816102828, "email from format error", nil)
	}
	service_pass := c.GetString("service_pass")
	if(len(service_pass)<1){
		c.ErrorJson(20220815100423, "get email pass error", nil)
	}
	service_host := c.GetString("service_host")
	if(len(service_host)<1){
		c.ErrorJson(20220815100527,"get email host error",nil)
	}
	service_port := c.GetString("service_port")
	if(len(service_port)<1){
		c.ErrorJson(20220815100631,"get email servive port error",nil)
	}
	campaign_id,camErr := c.GetInt64("campaign_id",0)
	if(camErr!=nil){
		c.ErrorJson(202208151003, "get campaign id error", nil)
	}
	// sender := c.GetString("sender_name")
	// if(len(sender)<1){
	// 	c.ErrorJson(202208151003, "get sender name error", nil)
	// }
	var emailCampaign =&models.Campaign{}
	var ecErr error
	if(campaign_id!=0){
	///vail campaign id valid
	camModel:=models.Campaign{}
	emailCampaign,ecErr=camModel.FindCambyid(campaign_id)
	if(ecErr!=nil){
		c.ErrorJson(20220815101442, "can not find the campaign by the id", nil)
	}
	}
	emailSer:=models.EmailService{
		Name: service_name,
		From: service_from,
		Password: service_pass,
		Port: service_port,
		Host: service_host,
		// SenderName: sender,
		Campaign: emailCampaign,
		Status: 1,
		AccountId:  &models.Account{Id: accountId},
	}
	emId, emErr:=emailSer.Createemailser(emailSer)
	if(emErr!=nil){
		c.ErrorJson(20220815101754, emErr.Error(), nil)
	}
	c.SuccessJson(dto.IdResponse{Id: emId})
}
//get email service by id

func (c *EmailserviceController) Getemailservice(){
	id, _ := c.GetInt64(":id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailSer := models.EmailService{}
	serEntity, err := emailSer.GetEmailServiceById(id, accountId)
	if err != nil {
		c.ErrorJson(20220815102320, "get email service error", nil)
	}
	emailServiceEntityDto:=dto.EmailServiceEntityDto{
		Id: serEntity.Id,
		From: serEntity.From,
		Password: serEntity.Password,
		Host: serEntity.Host,
		Port: serEntity.Port,
		Name: serEntity.Name,
		Ssl: serEntity.Ssl,
	}
	c.SuccessJson(emailServiceEntityDto)
}
