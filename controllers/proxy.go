package controllers

import (
	// "marketing/models"
	"encoding/json"
	"marketing/dto"
	"marketing/models"

	// "github.com/beego/beego/v2/core/utils"
	"marketing/utils"

	// "github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/i18n"
)

type ProxyController struct {
	BaseController
	i18n.Locale
}


func (c *ProxyController) ChildPrepare() {

}

//get proxy list
func (c *ProxyController) GetProxyList() {
	var proxy models.Proxy
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	start, _ := c.GetInt("page", 0)
	size, _ := c.GetInt("size", 10)
	search := c.GetString("search")
	// var proxylistresp ProxyListresp
	// proxylistresp.Total = 0
	// proxylistresp.Records = make([]models.Proxy, 0)
	proxylist, err := proxy.Getproxybyaccount(accountId, start, size, search)
	if err != nil {
		c.ErrorJson(202108031641149, err.Error(), nil)
	}
	// get proxy count by account
	total, err := proxy.GetProxyCountbyaccount(accountId, search)
	var prodtolist []dto.ProxyDto
	for _, v := range proxylist {
		prodtolist = append(prodtolist, dto.ProxyDto{
			Id:          v.Id,
			Host:        v.Host,
			Port:        v.Port,
			User:        v.User,
			Pass:        v.Pass,
			Protocol:    v.Protocol,
			CountryCode: v.CountryCode,
			Addtime: v.Addtime.Format("2006-01-02 15:04:05"),
		})
	}
	proxylistresp := dto.ProxyRespDto{
		Total:   total,
		Records: prodtolist,
	}

	c.SuccessJson(proxylistresp)
}

//delete proxy api
func (c *ProxyController) DeleteProxy() {
	var proxy models.Proxy
	id, _ := c.GetInt64("id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	err := proxy.DeleteProxy(id, accountId)
	if err != nil {
		c.ErrorJson(202403191109, err.Error(), nil)
	}
	c.SuccessJson(nil)
}

//get proxy detail
func (c *ProxyController) GetProxyDetail() {
	var proxy models.Proxy
	id, _ := c.GetInt64("id")
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	proxydetail, err := proxy.GetProxyDetail(id, accountId)
	if err != nil {
		c.ErrorJson(202403191109, err.Error(), nil)
	}
	proxyDetaildto := dto.ProxyDto{
		Id:          proxydetail.Id,
		Host:        proxydetail.Host,
		Port:        proxydetail.Port,
		User:        proxydetail.User,
		Pass:        proxydetail.Pass,
		Protocol:    proxydetail.Protocol,
		CountryCode: proxydetail.CountryCode,
	}
	c.SuccessJson(proxyDetaildto)
}

var validProtol = []string{"http", "https", "socket5"}

//save proxy api
func (c *ProxyController) SaveProxy() {
	var proxy models.Proxy
	var proxydto dto.ProxyDto
	err := c.ParseForm(&proxydto)
	if err != nil {
		c.ErrorJson(202403191109, err.Error(), nil)
	}
	if !utils.Contains(validProtol, proxydto.Protocol) {
		c.ErrorJson(202403201002102, "Invalid protocol", nil)
		return
	}

	uid := c.GetSession("uid")
	accountId := uid.(int64)
	proxy.Id = proxydto.Id
	proxy.Host = proxydto.Host
	proxy.Port = proxydto.Port
	proxy.User = proxydto.User
	proxy.Pass = proxydto.Pass
	proxy.Protocol = proxydto.Protocol
	proxy.CountryCode = proxydto.CountryCode
	proxy.Account = &models.Account{Id: accountId}
	if proxy.Id > 0 {
		pn, perr := proxy.GetProxyDetail(proxy.Id, accountId)
		if perr != nil {
			c.ErrorJson(202403200947, perr.Error(), nil)
		}
		if pn.Account.Id != accountId {
			c.ErrorJson(202403200948115, "no permission to update other's proxy", nil)
		}
		err = proxy.UpdateProxy(proxy)
		if err != nil {
			c.ErrorJson(202403200945112, err.Error(), nil)
		}
	} else {
		nid, err := proxy.Save(proxy)
		if err != nil {
			c.ErrorJson(202403191109, err.Error(), nil)
		}
		proxy.Id = nid
	}
	saveProxyDto := dto.SaveProxyDto{
		Id: proxy.Id,
	}
	c.SuccessJson(saveProxyDto)
}

//get protocol list
func (c *ProxyController) GetProtollist() {
	c.SuccessJson(validProtol)
}

//import proxy list json
func (c *ProxyController) ImportProxyList() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	proxyList := []dto.ImportProxyDto{}
	// logs.Info(c.Ctx.Input.RequestBody)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &proxyList); err == nil {
		for _, v := range proxyList {

			// v.Account = &models.Account{Id: accountId}
			if len(v.Host) <= 0 || len(v.Port) <= 0 || len(v.Protocol) <= 0 {
				c.ErrorJson(202403291411154, "Invalid proxy data", nil)
			}
			proxyItem:=models.Proxy{
				Host:v.Host,
				Port:v.Port,
				User:v.User,
				Pass:v.Pass,
				Protocol:v.Protocol,
				Account: &models.Account{Id: accountId},
			}

			
			_, errs := proxyItem.Save(proxyItem)
			if errs != nil {
				c.ErrorJson(202403191109, err.Error(), nil)
			}
		}
		c.SuccessJson(true)
	} else {
		
		c.ErrorJson(202403291414163, err.Error(), nil)
	}

}
//count proxy by account id
func (c *ProxyController) CountProxyByAccount() {
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	var proxy models.Proxy
	total, err := proxy.GetProxyCountbyaccount(accountId, "")
	if(err!=nil){
		c.ErrorJson(202408211352195, err.Error(), nil)
	}
	proxyCount:=dto.ProxyCountDto{
		Total: total,
	}
	c.SuccessJson(proxyCount)

}
