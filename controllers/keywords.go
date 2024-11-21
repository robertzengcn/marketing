package controllers

import (
	// "marketing/models"	
	"github.com/beego/i18n"
	"marketing/models"
	"strconv"
)

type KeywordsController struct {
	BaseController
	i18n.Locale
}

//get keyword list
func (c *KeywordsController) ListKeywordsbytag() {
	number := c.GetString("number")
	num, nerr := strconv.Atoi(number)
	if(nerr!=nil){
		c.ErrorJson(202402131114,nerr.Error(),nil)
	}
	uid := c.GetSession("uid")
	accountId:=uid.(int64)
	keywordModel:=models.Keyword{}
	tags := []string{}
	//Getkeywordbytag		
	inputValues, _ := c.Input()
	for k, v := range inputValues {
		if k == "tags[]" {
			if len(v) > 0 {
				tags = append(tags, v...)
			}
		} 		
	}
	// Logs.Info(tags)
	if(len(tags)<=0){
		c.ErrorJson(20240213115036,"tag number not enough", nil)
	}
	keywordArr,kerr:=keywordModel.Getkeywordbytag(tags,num,accountId)
	if(kerr!=nil){
		c.ErrorJson(20240213105532,kerr.Error(), nil)
	}
	keywordsArr:=[]string{}
	for _,v:=range keywordArr{
		keywordsArr=append(keywordsArr,v.Keyword)
	}

	c.SuccessJson(keywordsArr)
}