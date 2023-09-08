package controllers

import (
	"marketing/models"	
	"github.com/beego/i18n"
	"marketing/utils"
	// "github.com/beego/beego/v2/core/logs"
	// "fmt"
)
type ScraperinfoController struct {
	BaseController
	i18n.Locale
}

//save scraper info api
func (c *ScraperinfoController) SaveScraperInfo(){
	types := c.GetString("title")
	content := c.GetString("content")
	url := c.GetString("url")
	if(len(url)<=2){
		c.ErrorJson(20220619100351,"url not exist",nil)
	}
	//check url whether start with //
	if(url[0:2]=="//"){
		url="https:"+url
	}
	lang := c.GetString("lang","zh-cn")
	allowLang:=[]string{"en-us","zh-cn"}
	if(!utils.Contains(allowLang,lang)){
		c.ErrorJson(20220621100351,"lang not exist",nil)
	}
	socialtask_id,_ := c.GetInt64("socialtask_id")
	if(socialtask_id>0){
		//valid social task id
		socialtaskModel:=models.SocialTask{}
		_,serr:=socialtaskModel.GetSocialTaskById(socialtask_id)
		if(serr!=nil){
			c.ErrorJson(202306200901,serr.Error(),nil)
		}
	}
	scraperinfoModel:=models.ScraperInfo{}
	scraperinfoVar:=models.ScraperInfo{Title:types,Content:content,Url:url,Lang:lang,SocialTaskId:&models.SocialTask{Id:socialtask_id}}
	sid,serr:=scraperinfoModel.SavedataDb(&scraperinfoVar)
	if(serr!=nil){
		c.ErrorJson(20220619100953,serr.Error(),nil)
	}
	// arr := []int{1, 2, 3, 4}
	// for i := range arr {
	// 	arr = append(arr, i)
	// }
	// fmt.Println(arr)
	c.SuccessJson(sid)
}
func (c *ScraperinfoController) Getscrapyinfolist(){
	sotaskid,_ := c.GetInt64("sotaskid",0)

	if(sotaskid<=0){
		c.ErrorJson(20290705093655,"social task id emtpy",nil)
	}
	limitNum,_ := c.GetInt("limit",5)

	scraperinfoModel:=models.ScraperInfo{}
	
	scralist,scerr:=scraperinfoModel.GetScraperInfoByTaskId(sotaskid,limitNum,true,"usetime asc")
	if(scerr!=nil){
		c.ErrorJson(20230705095362,scerr.Error(),nil)
	}
	c.SuccessJson(scralist)
}
//update the scraper info process time
func (c *ScraperinfoController) Updatescrapyinfoprocess(){
	id,_ := c.GetInt64("id",0)
	if(id<=0){
		c.ErrorJson(202309070955,"id emtpy",nil)
	}
	scraperinfoModel:=models.ScraperInfo{}
	_,serr:=scraperinfoModel.FindOne(id)
	if(serr!=nil){
		c.ErrorJson(202309070959,serr.Error(),nil)
	}
	scraperinfoModel.UpdateProcesstime(id)
	c.SuccessJson(nil)
}
