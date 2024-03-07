package routers

import (
	"fmt"
	"marketing/controllers"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

func init() {
	langs, err := beego.AppConfig.String("langs::types") // 1
	if err != nil {                                      // 2
		fmt.Println("Failed to load languages from the app.conf")
		// return
	}
	// fmt.Println("11111")
	langsArr := strings.Split(langs, "|")
	// fmt.Printf("%v", langsArr)
	for _, lang := range langsArr {
		// fmt.Print(lang)
		if err := i18n.SetMessage(lang, "conf/"+lang+".ini"); err != nil { // 5
			fmt.Println("Failed to set message file for l10n")
			return
		}
	}
	beego.Router("*", &controllers.MainController{}, "options:Checkoption")

	beego.Router("/", &controllers.MainController{})

	admin := beego.NewNamespace("/admin",
		beego.NSInclude(
			&controllers.TestController{},
		),
		beego.NSRouter("/campaign/create", &controllers.CampaignController{}, "post:CreateCampaign"),
		beego.NSRouter("/task/create", &controllers.TaskController{}, "post:CreateTask"),
		beego.NSRouter("/login/accountlogin", &controllers.AccountController{}, "post:Validaccount"),
		beego.NSRouter("/task/updatetask", &controllers.TaskController{}, "post:UpdateTaskstatus"),
		//list campagin
		beego.NSRouter("/campaign", &controllers.CampaignController{}, "get:ListCampaign"),
		// beego.Router("/addSite", &controllers.CampaignController{}, "post:Createsite")
		beego.NSRouter("/emailtpl/create", &controllers.EmailtplController{}, "post:CreateEmailtpl"),
		beego.NSRouter("/emailservice/add", &controllers.EmailserviceController{}, "post:Addemailservice"),
		beego.NSRouter("/schedule/add", &controllers.ScheduleController{}, "post:CreateSchedule"),

		// beego.Router("/welcome", &controllers.CampaignController{}, "get:Welcome")
		beego.NSRouter("/healthcheck", &controllers.MainController{}, "get:Healthcheck"),
		// beego.Router("/social/create", &controllers.SocialController{}, "post:CreateSocialAccount")

		beego.NSRouter("/schedule/list", &controllers.ScheduleController{}, "get:ListSchedule"),
		// beego.NSRouter("/socialcampaign/create", &controllers.CampaignController{}, "post:CreateSocialAccount"),
		// beego.NSRouter("/getsobyCam", &controllers.CampaignController{}, "get:GetSocialAccount"),

		//beego.Router("/getsobyCam", &controllers.CampaignController{}, "get:GetSocialAccount")
	)
	beego.AddNamespace(admin)
	ns := beego.NewNamespace("/test",
		beego.NSInclude(
			&controllers.TestController{},
		),
		beego.NSRouter("/savesearchreq", &controllers.TestController{}, "post:Savesearchrequest"),
		beego.NSRouter("/testsendemail", &controllers.TestController{}, "post:Testsendemail"),
		beego.NSRouter("/checkemailsend", &controllers.TestController{}, "post:Checkemailsend"),
		beego.NSRouter("/testgetadultkeyword", &controllers.TestController{}, "post:Getadultkeyword"),
		beego.NSRouter("/createtastschedule", &controllers.TestController{}, "post:Createtasksched"),
		beego.NSRouter("/createdaytask", &controllers.TestController{}, "post:CreatedayTask"),
		beego.NSRouter("/getkeywordbytag", &controllers.TestController{}, "get:Getkeywordbytag"),
		beego.NSRouter("/getkeywordapi", &controllers.TestController{}, "get:Getkeywordapi"),
		beego.NSRouter("/loademailapi", &controllers.TestController{}, "post:LoadEmailapi"),
		beego.NSRouter("/importkeyword", &controllers.TestController{}, "post:Loadkeywordapi"),
		beego.NSRouter("/getproxylist", &controllers.TestController{}, "get:TestProxylist"),
		beego.NSRouter("/updateproxy", &controllers.TestController{}, "get:UpdatemulProxy"),
		beego.NSRouter("/getemailser", &controllers.TestController{}, "post:Getemailbycampaign"),
		beego.NSRouter("/changeproxy", &controllers.TestController{}, "get:ChangeProxy"),
		beego.NSRouter("/handleproxy", &controllers.TestController{}, "get:GetProxylist"),
		beego.NSRouter("/campaign", &controllers.CampaignController{}, "get:ListCampaign"),
	)
	beego.AddNamespace(ns)
	api := beego.NewNamespace("/api",
		beego.NSInclude(
			&controllers.CampaignController{},
		),
		//get socoial account by campaign id
		// beego.NSRouter("/getsobyCam", &controllers.CampaignController{}, "get:GetSocialAccount"),
		beego.NSRouter("/listsoCampaign", &controllers.SocialController{}, "get:Listsocialcampaigin"),
		beego.NSRouter("/listsotask", &controllers.SocialController{}, "get:Getsocialtasklist"),
		beego.NSRouter("/savesolink", &controllers.ScraperinfoController{}, "post:SaveScraperInfo"),
		beego.NSRouter("/getscrapeinfolist", &controllers.ScraperinfoController{}, "get:Getscrapyinfolist"),
		beego.NSRouter("/getsocialtaskinfo", &controllers.SocialController{}, "get:Getsocialtaskinfo"),
		beego.NSRouter("/taskkeyword", &controllers.SocialController{}, "get:Gettaskkeyword"),
		beego.NSRouter("/savesocialtask", &controllers.SocialController{}, "post:Savesocialtask"),
		beego.NSRouter("/updatescrapeprotime", &controllers.ScraperinfoController{}, "post:Updatescrapyinfoprocess"),
		// beego.NSRouter("/login", &controllers.AccountController{}, "post:Validaccount"),
		beego.NSRouter("/campaign", &controllers.CampaignController{}, "get:ListCampaign"),
		beego.NSRouter("/user/info", &controllers.AccountController{}, "get:Accountinfo"),
		//list social task type api
		beego.NSRouter("/socialtasktype", &controllers.SocialController{}, "get:Listsocialtasktype"),
		beego.NSRouter("/tag", &controllers.TagController{}, "get:ListTag"),
		beego.NSRouter("/getkeywordtag", &controllers.KeywordsController{}, "post:ListKeywordsbytag"),
		beego.NSRouter("/user/signout", &controllers.AccountController{}, "get:Signout"),
		beego.NSRouter("/socialaccount", &controllers.SocialAccountController{}, "post:Savesocialaccount"),
		beego.NSRouter("/socialaccount", &controllers.SocialAccountController{}, "get:GetSocialAccount"),
		beego.NSRouter("/socialaccount/list", &controllers.SocialAccountController{}, "get:Listsocialaccount"),
		beego.NSRouter("/socialaccount/update", &controllers.SocialAccountController{}, "post:Updatesocialaccount"),
	)
	beego.AddNamespace(api)

	user := beego.NewNamespace("/user",
		beego.NSInclude(
			&controllers.TestController{},
		),
		beego.NSRouter("/login", &controllers.AccountController{}, "post:Validaccount"),
	)
	beego.AddNamespace(user)

}
