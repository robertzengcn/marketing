package controllers
import (
	"github.com/beego/i18n"
	"marketing/models"
	"marketing/utils"
	"marketing/dto"
	"strings"
)

type EmailserviceController struct {
	BaseController
	i18n.Locale
}

func (c *EmailserviceController) ChildPrepare(){

}

///add email service
func (c *EmailserviceController) Addemailservice(){
	id, _ := c.GetInt64(":id",0)
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	service_name := c.GetString("service_name")
	if(len(service_name)<1){
		c.ErrorJson(20241024140626, "get service name error", nil)
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
	if(len(service_port)>5){
		c.ErrorJson(20241024142548,"port length is too long",nil)
	}
	ssl,sslErr := c.GetInt8("service_ssl",0)
	if(sslErr!=nil){
		c.ErrorJson(20241022143549, "get ssl error", nil)
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
	upass,uerr:=utils.Encrypt(service_pass)

	if(uerr!=nil){
		c.ErrorJson(20241022102565, uerr.Error(), nil)
	}
	emailSer:=models.EmailService{
		Name: service_name,
		From: service_from,
		Password: upass,
		Port: service_port,
		Host: service_host,
		// SenderName: sender,
		Campaign: emailCampaign,
		Status: 1,
		AccountId:  &models.Account{Id: accountId},
		Ssl: ssl,
	}
	var emId int64
	var emErr error
	if(id!=0){
		emId= id
		emailSer.Id=id
		emErr=emailSer.UpdateEmailService(&emailSer)
		if(emErr!=nil){
			c.ErrorJson(20241022110187, emErr.Error(), nil)
		}
	}else{
	emId, emErr=emailSer.Createemailser(emailSer)
	}
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
		c.ErrorJson(202410241006112, "get email service error", nil)
	}
	//decrypted password
	passwrod,uerr:=utils.Decrypt(serEntity.Password)
	if(uerr!=nil){
		c.ErrorJson(20241022102798, uerr.Error(), nil)
	}
	emailServiceEntityDto:=dto.EmailServiceEntityDto{
		Id: serEntity.Id,
		From: serEntity.From,
		Password: passwrod,
		Host: serEntity.Host,
		Port: serEntity.Port,
		Name: serEntity.Name,
		Ssl: serEntity.Ssl,
	}
	c.SuccessJson(emailServiceEntityDto)
}
//delete email service
func (c *EmailserviceController) Deleteemailservice(){
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	id, terr := c.GetInt64(":id")
	if terr != nil {
		c.ErrorJson(20220618162124, terr.Error(), nil)
	}
	emailSer := models.EmailService{}
	eres, ereserr := emailSer.GetEmailServiceById(id, accountId)
	if ereserr != nil {
		c.ErrorJson(20220815102320, ereserr.Error(), nil)
	}
	if eres == nil {
		c.ErrorJson(202410241406144, "email service not exist", nil)
	}
	err := emailSer.DeleteEmailService(id,accountId)
	if err != nil {
		c.ErrorJson(202410241407148, err.Error(), nil)
	}
	
	c.SuccessJson(dto.IdResponse{Id: id})
}

///get email service list by account id
func (c *EmailserviceController) GetServiceList() {
	page, perr := c.GetInt64("page", 0)
	if perr != nil {
		c.ErrorJson(202410221123152, perr.Error(), nil)
	}
	size, serr := c.GetInt64("size", 10)
	if serr != nil {
		c.ErrorJson(202410221123156, serr.Error(), nil)
	}
	search:=c.GetString("search","")

	orderby := c.GetString("orderby","")
	neworderby := strings.ReplaceAll(orderby, "-", "")
	if len(neworderby) > 0 {

		orderbyvaild := utils.Contains([]string{"id"}, neworderby)
		if !orderbyvaild {
			c.ErrorJson(202410221124166, "orderby incorrect", nil)
		}
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailserModel := models.EmailService{}
	emailtplList, emailerr := emailserModel.GetEmailServiceListByAccountId(accountId, page, size,search,neworderby)
	if emailerr != nil {
		c.ErrorJson(20220617155924, emailerr.Error(), nil)
	}
	var emailtpldtos []dto.EmailServiceListDto
	for _, s := range emailtplList {
		emailtpldtos = append(emailtpldtos, dto.EmailServiceListDto{
			Id:      s.Id,
			Name:   s.Name,
			From: s.From,
			Host: s.Host,
			CreateTime: s.Created.Format("2006-01-02 15:04:05"),		
		})
	}
	ecnum,ecerr:=emailserModel.CountEmailService(accountId,search)
	if ecerr != nil {
		c.ErrorJson(202409271129118, ecerr.Error(), nil)
	}

	// Removed unused variable resp and fixed instantiation issue
	resp := dto.CommonResponse[[]dto.EmailServiceListDto]{
		Record: emailtpldtos,
		Total:    ecnum,
	}

	c.SuccessJson(resp)
	// c.SuccessJson(emailtpldtos)
}

