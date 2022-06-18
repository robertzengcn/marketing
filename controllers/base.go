package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	// "github.com/beego/beego/v2/core/logs"
	// "reflect"
	// "fmt"
	"github.com/beego/i18n"
	"marketing/utils"
	"fmt"
	"strings"
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
	Status bool
	Code   int
	Msg    string
	Data   interface{}
}

type langType struct {
	Lang string
	Name string
}

func (c *BaseController) SuccessJson(data interface{}) {

	res := ReturnMsg{
		true, 200, "success", data,
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

func Filter_user(ctx *context.Context) {
	s := []string{"/login", "/login/accountlogin","/healthcheck"} //defined url that not need to valid user login
	// l := logs.GetLogger()
	// res:=ctx.Input.Session("uid")
	// l.Println(res)
	// fmt.Println(reflect.TypeOf(res))
	_, ok := ctx.Input.Session("uid").(int64)

	// l.Println(id)
	// l.Println(ok)
	if !ok { //user not login
		if !utils.Contains(s, ctx.Request.RequestURI) {
			// if ctx.Input.IsAjax() {
				jsonData := make(map[string]interface{}, 3)

				jsonData["errcode"] = 403
				jsonData["message"] = "You have to login to continue"

				returnJSON, _ := json.Marshal(jsonData)
				ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8;")
				ctx.ResponseWriter.Write(returnJSON)
			// } else {
			// 	ctx.Redirect(302, "/login")
			// }
		}
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
