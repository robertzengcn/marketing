package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
	"marketing/dto"
	"github.com/beego/beego/v2/core/logs"
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
	emailFilter := models.EmailFilter{}
	emailFilter.Name = c.GetString("name")
	emailFilter.AccountId = &models.Account{Id: accountId}
	id,err := emailFilter.CreateEmailFilter(&emailFilter)
	if err != nil {
		c.ErrorJson(20220815102320, "create email filter error", nil)
	}
	//config filter detail list
	var filterdetails []string
		//update account proxy
	inputValues, _ := c.Input()
	for k, v := range inputValues {
		if k == "filterdetails[]" {
			if len(v) > 0 {
				// Convert v from []string to []int64
				filterdetails = append(filterdetails, v...)
			}
		}
	}
	logs.Info(filterdetails)
	c.SuccessJson(dto.IdResponse{Id: id})
}
//get email filter by id
func (c *EmailFilterController) GetEmailFilterById() {
	id, _ := c.GetInt64("id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	emailFilter := models.EmailFilter{}
	FilterEntity, err := emailFilter.GetEmailFilterById(id,accountId)
	if err != nil {
		c.ErrorJson(20220815102320, "get email filter error", nil)
	}
	c.SuccessJson(FilterEntity)
}
