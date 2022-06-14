package controllers

import (
	"marketing/models"
	"github.com/beego/i18n"
	"path/filepath"
	"runtime"
	"fmt"
)

type TestController struct {
	BaseController
	i18n.Locale
}
///test save search request
func(c *TestController) Savesearchrequest(){
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	filename:=apppath+"/tests/20220509-threaded-results.json"
	res,rerr:=models.DefaultTask.Readfile(filename)
	if(rerr!=nil){
		panic(rerr.Error())	
	}
	searchreqModel:=models.SearchRequest{}
	serr:=searchreqModel.Savesrlist(res)
	if(serr!=nil){
		fmt.Println(serr)
	}
	c.Ctx.WriteString("success")
}
func(c *TestController) Callemailscrape(){
	
}
