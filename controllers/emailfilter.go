package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"

	"encoding/json"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/logs"
	"marketing/dto"
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