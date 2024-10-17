package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"strings"
	"encoding/json"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/logs"
	"marketing/dto"
	"marketing/utils"
)

type EmailFilterController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

//create email filter
func (c *EmailFilterController) CreateEmailFilter() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	// emailFilter := models.EmailFilter{}
	// emailFilter.Name = c.GetString("name")
	// emailFilter.AccountId = &models.Account{Id: accountId}
	// id,err := emailFilter.CreateEmailFilter(&emailFilter)
	// if err != nil {
	// 	c.ErrorJson(20220815102320, "create email filter error", nil)
	// }
	// //config filter detail list
	// var filterdetails []string
	// 	//update account proxy
	// inputValues, _ := c.Input()
	// for k, v := range inputValues {
	// 	if k == "filterdetails[]" {
	// 		if len(v) > 0 {
	// 			// Convert v from []string to []int64
	// 			filterdetails = append(filterdetails, v...)
	// 		}
	// 	}
	// }
	// logs.Info(filterdetails)
	// //save filter detail
	// emailFilterDetail := models.EmailFilterDetail{}
	// for _, v := range filterdetails {
	// 	emailFilterDetail.FilterId = &models.EmailFilter{Id: id}
	// 	emailFilterDetail.AccountId = &models.Account{Id: accountId}
	// 	emailFilterDetail.Content = v
	// 	emailFilterDetail.CreateEmailFilterDetail(&emailFilterDetail)
	// }
	// c.SuccessJson(dto.IdResponse{Id: id})
	emailFilterdto := dto.EmailFilterEntityDto{}

	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &emailFilterdto); err == nil {
		emailFilter := models.EmailFilter{}
		emailFilter.Name = emailFilterdto.Name
		emailFilter.Description=emailFilterdto.Description
		emailFilter.AccountId = &models.Account{Id: accountId}
		id, err := emailFilter.CreateEmailFilter(&emailFilter)
		if err != nil {
			c.ErrorJson(20241014144062, "create email filter error", nil)
		}
		for _, v := range emailFilterdto.FilterDetails {

			emailFilterDetail := models.EmailFilterDetail{}
			
			emailFilterDetail.FilterId = &models.EmailFilter{Id: id}
			emailFilterDetail.AccountId = &models.Account{Id: accountId}
			emailFilterDetail.Content = v.Content
			if(v.Id==0){
				
			
				eId,eferr:=emailFilterDetail.CreateEmailFilterDetail(&emailFilterDetail)
				if(eferr!=nil){
					c.ErrorJson(20241014144276, eferr.Error(), nil)
				}else{
					emailFilterDetail.Id=eId
				}
			}else{
				
				emailFilterDetail.Id=v.Id

				emailFilterDetail.UpdateEmailFilterDetail(&emailFilterDetail)
			}
		}
		c.SuccessJson(dto.IdResponse{Id: id})
	}else{
		logs.Error(err)
		c.ErrorJson(20241014144190, err.Error(), nil)
	}
	

}

//get email filter by id
func (c *EmailFilterController) GetEmailFilterById() {

	id, _ := c.GetInt64(":id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailFilter := models.EmailFilter{}
	FilterEntity, err := emailFilter.GetEmailFilterById(id, accountId)
	if err != nil {
		c.ErrorJson(20220815102320, "get email filter error", nil)
	}
	filterDetailModel := models.EmailFilterDetail{}

	fitarr, _ := filterDetailModel.GetEmailFilterDetailByFilterId(id, accountId)
	efdto := dto.EmailFilterEntityDto{
		Id:   FilterEntity.Id,
		Name: FilterEntity.Name,
		Description: FilterEntity.Description,
	}
	// if(fitarr!=nil){
	for _, v := range fitarr {
		efdto.FilterDetails = append(efdto.FilterDetails, dto.EmailFilterDetailDto{
			Id:      v.Id,
			Content: v.Content,
		})
	}
	// }
	c.SuccessJson(efdto)
}
//update email filter by id
func (c *EmailFilterController) UpdateEmailFilter() {
	id, _ := c.GetInt64(":id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailFilterdto := dto.EmailFilterEntityDto{}

	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &emailFilterdto); err == nil {
		emailFilter := models.EmailFilter{}
		emailFilter.Id = id
		emailFilter.Name = emailFilterdto.Name
		emailFilter.Description=emailFilterdto.Description
		emailFilter.AccountId = &models.Account{Id: accountId}
		err := emailFilter.UpdateEmailFilter(&emailFilter)
		if err != nil {
			c.ErrorJson(202410141451139, err.Error(), nil)
		}
		for _, v := range emailFilterdto.FilterDetails {

			emailFilterDetail := models.EmailFilterDetail{}
		
			emailFilterDetail.FilterId = &models.EmailFilter{Id: id}
			emailFilterDetail.AccountId = &models.Account{Id: accountId}
			emailFilterDetail.Content = v.Content
			if(v.Id==0){
				
			
				eId,eferr:=emailFilterDetail.CreateEmailFilterDetail(&emailFilterDetail)
				if(eferr!=nil){
					c.ErrorJson(20241014144276, eferr.Error(), nil)
				}else{
					emailFilterDetail.Id=eId
				}
			}else{
				
				emailFilterDetail.Id=v.Id

				emailFilterDetail.UpdateEmailFilterDetail(&emailFilterDetail)
			}
		}
		c.SuccessJson(dto.IdResponse{Id: id})
	}
}
//delete email filter by id
func (c *EmailFilterController) DeleteEmailFilter() {
	id, _ := c.GetInt64(":id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailFilterDetailModel:=models.EmailFilterDetail{}
	//delete email fileter detail first
	ederr:=emailFilterDetailModel.DeleteEmailFilterDetail(id,accountId)
	if(ederr!=nil){
		c.ErrorJson(202410141500, ederr.Error(), nil)
	}
	//delete email filter
	emailFilter := models.EmailFilter{}
	err := emailFilter.DeleteEmailFilter(id, accountId)
	if err != nil {
		c.ErrorJson(202410141452139, err.Error(), nil)
	}
	c.SuccessJson(dto.IdResponse{Id: id})
}
//list email filter
func (c *EmailFilterController) ListEmailFilter() {
	page, perr := c.GetInt64("page", 0)
	if perr != nil {
		c.ErrorJson(202410161508191, perr.Error(), nil)
	}
	size, serr := c.GetInt64("size", 10)
	if serr != nil {
		c.ErrorJson(202410161509195, serr.Error(), nil)
	}
	search:=c.GetString("search","")

	orderby := c.GetString("orderby","")
	neworderby := strings.ReplaceAll(orderby, "-", "")
	if len(neworderby) > 0 {

		orderbyvaild := utils.Contains([]string{"id"}, neworderby)
		if !orderbyvaild {
			c.ErrorJson(202410161509205, "orderby incorrect", nil)
		}
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailFilter := models.EmailFilter{}
	emailFilterList, err := emailFilter.ListEmailFilter(accountId, page, size, search, orderby)
	num,ecerr:=emailFilter.CountEmailFilter(accountId, search)
	if(ecerr!=nil){
		c.ErrorJson(202410171506217, ecerr.Error(), nil)
	}
	if err != nil {
		c.ErrorJson(202410161510213, err.Error(), nil)
	}
	efdto := []dto.EmailFilterEntityDto{}
	for _, v := range emailFilterList {
		efdto = append(efdto, dto.EmailFilterEntityDto{
			Id:   v.Id,
			Name: v.Name,
			CreatedTime: v.Created.Format("2006-01-02 15:04:05"),
		})
	}
	resp := dto.CommonResponse[[]dto.EmailFilterEntityDto]{
		Record: efdto,
		Total:    num,
	}
	c.SuccessJson(resp)
}