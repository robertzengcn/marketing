package controllers

import (
	"encoding/json"
	"errors"

	// "github.com/beego/beego/v2/core/logs"
	//  "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	// "github.com/beego/beego/v2/core/logs"
	// "reflect"
	// "fmt"
	"fmt"
	"marketing/models"
	"marketing/utils"
	"strings"

	"github.com/beego/i18n"
)

type Controllerreturn interface {
	SuccessJson()
	ErrorJson()
}

type BaseController struct {
	beego.Controller
	i18n.Locale
}

type ReturnMsg struct {
	Status bool `json:"status"`
	Code   int  `json:"code"`
	Msg    string `json:"msg"`
	Data   interface{} `json:"data"`
}

type langType struct {
	Lang string
	Name string
}

func (c *BaseController) SuccessJson(data interface{}) {

	res := ReturnMsg{
		true, 20000, "success", data,
	}
	c.Data["json"] = res
	c.ServeJSON() //对json进行序列化输出
	c.StopRun()
}

func (c *BaseController) ErrorJson(code int, msg string, data interface{}) {

	res := ReturnMsg{
		false, code, msg, data,
	}

	c.Data["json"] = res
	c.ServeJSON() //对json进行序列化输出
	c.StopRun()
}
///valid bearer token,return account
func Valid_beartoken(authstr string)(*models.AccountToken,error){
	
	splitToken := strings.Split(authstr, "Bearer ")
	if(len(splitToken)<2){
		return nil,errors.New("token is empty")
	}
	token:=splitToken[1]
	accountTokenModel:=models.AccountToken{}
	return accountTokenModel.CheckAccounttoken(token)
}
///get account by token
func GetAccountbytoken(token string)(models.Account,error){
	accountTokenModel:=models.AccountToken{}
	atoken,_:=accountTokenModel.CheckAccounttoken(token)
	return *atoken.Account,nil
}

func Filter_user(ctx *context.Context) {
	s := []string{"/api/login", 
	// "/admin/login/accountlogin",
	"/admin/healthcheck",
	"/api/user/login",
	// "api/getsobyCam",
	} //defined url that not need to valid user login
	// l := logs.GetLogger()
	//defined basic auth array

	// _, ok := ctx.Input.Session("uid").(int64)

	// l.Println(id)
	// l.Println(ok)
	// if !ok { //user not login
		if !utils.Contains(s, ctx.Request.RequestURI) {
			// logs.Info("the request url is: "+ctx.Request.RequestURI)
			//check bearer token
			authstr := ctx.Input.Header("Authorization")
			// logs.Info("authstr is: "+authstr) 
			// logs.Info(authstr)
			if(len(authstr)>0){
				acctoken,accerr:=Valid_beartoken(authstr)
				
				if(accerr==nil&&acctoken!=nil){
					// logs.Info("uid is"+fmt.Sprintf("%d", acctoken.Account.Id))
					ctx.Output.Session("uid", acctoken.Account.Id)
					ctx.Output.Session("tokenId", acctoken.TokenId)
					return
				}
			}
			//change http response to 403

				jsonData := make(map[string]interface{}, 3)
				
				jsonData["status"]=false
				jsonData["code"] = 403
				jsonData["msg"] = "You have to login to continue"

				returnJSON, _ := json.Marshal(jsonData)
				ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8;")
				ctx.ResponseWriter.Write(returnJSON)
				ctx.ResponseWriter.WriteHeader(403)
		}
	// }
}
//basic Authorization
func Filter_basic(ctx *context.Context) {
	username, password, ok :=ctx.Request.BasicAuth()
	//  logs.Info("basic auth check")
	// logs.Info(password)
	if ok {
		fres:=Filter_account(username, password)
		
		if(!fres){
			// logs.Error(fres)
			Forbidenreturn(ctx)
		}
		// logs.Info("ok")
		return
	}
	//get basic auth from header
	//  authstr := ctx.Input.Header("Authorization")
	// //  logs.Info(authstr)
	//  if(len(authstr)<=0){
	// 	Forbidenreturn(ctx)
	//  }
	//  fres:=Filter_account(authstr)
	//  if(!fres){
	// 	 Forbidenreturn(ctx)
	//  }
	Forbidenreturn(ctx)
}
//return forbiden if user not login
func Forbidenreturn(ctx *context.Context){
	ctx.ResponseWriter.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
	ctx.ResponseWriter.WriteHeader(401)
	ctx.ResponseWriter.Write([]byte("401 Unauthorized\n"))
   
}
//check account valid
func Filter_account(username string,pass string)(bool){
	// autharr := strings.Split(authstr, ":")
	// //check array length
	// if(len(autharr)!=2){
	// 	return false
	// }
	// logs.Info(autharr)
	apiId,_:=models.DefaultApiauth.GetApiAuth(username,pass)
	if(apiId>0){
		return true
	}else{
		return false
	}
}

type ChildPrepareHolder interface {
	ChildPrepare()
}

func (c *BaseController) Prepare() {
	langstring, err := beego.AppConfig.String("langs::types")
	if err != nil {  // 2
		fmt.Println("Failed to load languages from the app.conf")
		// return
	}
	langs := strings.Split(langstring,"|")
	langNamestring, err := beego.AppConfig.String("langs::names")
	if err != nil {  // 2
		fmt.Println("Failed to load languages name from the app.conf")
		// return
	}
	names := strings.Split(langNamestring,"|")
	
	langTypes := make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	// names,_ := strings.Split(beego.AppConfig.String("lang::names"), "|")
	
	// isNeedRedir := false
	hasCookie := false

	lang := c.GetString("lang")

	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		hasCookie = true
	} 
	if !i18n.IsExist(lang) {
		lang = ""
		// isNeedRedir = false
		hasCookie = false
	}

	if len(lang) == 0 {
		al := c.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}
	if len(lang) == 0 {
		lang = "en-US"
		// isNeedRedir = false
	}
	curLang := langType{
		Lang: lang,
	}

	if !hasCookie {
		c.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}
	// l := logs.GetLogger()
	c.Lang = lang
	// l.Println("first time load language")
	// l.Println(lang)
    c.Data["Lang"] = curLang.Lang
    c.Data["CurLang"] = curLang.Name
    c.Data["RestLangs"] = restLangs


	//      l.Println("111111")
	// l.Println(c.Data["langTemplateKey"])
	// c.Data["Lang"]=c.Data["langTemplateKey"]
	c.EnableXSRF = false
	if app, ok := c.AppController.(ChildPrepareHolder); ok { // 5
		app.ChildPrepare()
	}
	// return isNeedRedir
}

// func (c *BaseController) setLangVer() bool {
//     isNeedRedir := false
//     hasCookie := false

//     // 1. Check URL arguments.
//     lang := c.GetString("lang")

//     // 2. Get language information from cookies.
//     if len(lang) == 0 {
//         lang = c.Ctx.GetCookie("lang")
//         hasCookie = true
//     } else {
//         isNeedRedir = true
//     }

//     // Check again in case someone modify by purpose.
//     if !i18n.IsExist(lang) {
//         lang = ""
//         isNeedRedir = false
//         hasCookie = false
//     }

//     // 3. Get language information from 'Accept-Language'.
//     if len(lang) == 0 {
//         al := c.Ctx.Request.Header.Get("Accept-Language")
//         if len(al) > 4 {
//             al = al[:5] // Only compare first 5 letters.
//             if i18n.IsExist(al) {
//                 lang = al
//             }
//         }
//     }

//     // 4. Default language is English.
//     if len(lang) == 0 {
//         lang = "en-US"
//         isNeedRedir = false
//     }

//     curLang := langType{
//         Lang: lang,
//     }

//     // Save language information in cookies.
//     if !hasCookie {
//         c.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
//     }

//     restLangs := make([]*langType, 0, len(langTypes)-1)
//     for _, v := range langTypes {
//         if lang != v.Lang {
//             restLangs = append(restLangs, v)
//         } else {
//             curLang.Name = v.Name
//         }
//     }

//     // Set language properties.
//     c.Lang = lang
//     c.Data["Lang"] = curLang.Lang
//     c.Data["CurLang"] = curLang.Name
//     c.Data["RestLangs"] = restLangs

//     return isNeedRedir
// }
