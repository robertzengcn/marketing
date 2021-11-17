package controllers

import(
	beego "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/server/web/context"
    "encoding/json"
    "github.com/beego/beego/v2/core/logs"
    "reflect"
    "fmt"
)

type Controllerreturn interface{
    SuccessJson()
    ErrorJson()
}

type BaseController struct {
    beego.Controller
}

type ReturnMsg struct {
	Status bool
    Code int
    Msg  string
    Data interface{}
}

func (this *BaseController) SuccessJson(data interface{}) {

    res := ReturnMsg{
        true,200, "success", data,
    }
    this.Data["json"] = res
    this.ServeJSON() //对json进行序列化输出
    this.StopRun()
}

func (this *BaseController) ErrorJson(code int, msg string, data interface{}) {

    res := ReturnMsg{
        false,code, msg, data,
    }

    this.Data["json"] = res
    this.ServeJSON() //对json进行序列化输出
    this.StopRun()
}

func Filter_user(ctx *context.Context) {
    l := logs.GetLogger()
    res:=ctx.Input.Session("uid")
    l.Println(res)
    fmt.Println(reflect.TypeOf(res))
    _, ok := ctx.Input.Session("uid").(int64)
    
    // l.Println(id)
    l.Println(ok)
    if !ok && ctx.Request.RequestURI != "/login" {
        if ctx.Input.IsAjax() {
            jsonData := make(map[string]interface{}, 3)

            jsonData["errcode"] = 403
            jsonData["message"] = "请登录后再操作"

            returnJSON, _ := json.Marshal(jsonData)

            ctx.ResponseWriter.Write(returnJSON)
        } else {
            ctx.Redirect(302, "/login")
        }
    }
}