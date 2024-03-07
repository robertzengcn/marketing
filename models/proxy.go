package models

import (
	// "encoding/base64"
	"strings"
	"time"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/parnurzeal/gorequest"
)

type Proxy struct {
	Id          int64     `orm:"pk;auto"`
	Host        string    `orm:"size(350)" valid:"Required"`
	Port        string    `orm:"size(6)" valid:"Required"`
	User        string    `orm:"size(350)"`
	Pass        string    `orm:"size(350)"`
	Protocol    string    `orm:"size(10)"`
	Available   int       `orm:"size(1);default(1);description(this mean status of the proxy)"`
	CountryCode string    `orm:"size(20)" json:"country_code"`
	Addtime     time.Time `orm:"auto_now_add;type(datetime)"`
	Checktime   time.Time `orm:"null;type(datetime)"`
	Source 		string    `orm:"size(20)"`
	Usetime	 time.Time `orm:"null;type(datetime)"`
	//Googleenable int  	 `orm:"size(1);default(0);description(whether the proxy is enable on google)"`
}

///defined table name
func (u *Proxy) TableName() string {
	return "proxy"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(Proxy))
}

// set engineer as INNODB
func (u *Proxy) TableEngine() string {
	return "INNODB"
}
//get proxy type
func (u *Proxy) Getproxytype() (Proxyway) {
	webproxy:=ProxyWebshare{}
	return &webproxy
}

//save proxy to database
func (u *Proxy) Save(proxy Proxy) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var proxyitem Proxy
	err := qs.Filter("host", proxy.Host).Filter("port", proxy.Port).Filter("user", proxy.User).Filter("pass", proxy.Pass).Filter("protocol", proxy.Protocol).One(&proxyitem)
	// logs.Error(err)
	if err == orm.ErrNoRows {
		id, err := o.Insert(&proxy)
		return id, err
	}
	return 0, err
}

//handle proxy list
func (u *Proxy) GetProxylist(pxw Proxyway) ([]Proxy, error) {
	return pxw.Proxylist()
}

//handle proxy from third party
func (u *Proxy) Handleproxy() ( error) {
	// pxw := ProxyWebshare{}
	//pxw := Asocksproxy{}
	pxw :=u.Getproxytype()
	proarr,perr:=u.GetProxylist(pxw)
	logs.Info(proarr)
	if(perr!=nil){
		logs.Error(perr)
		return nil
	}
	//range proxy list,save to database
	for _, proxy := range proarr {
		var proxyStr=proxy.Protocol+"://"+proxy.User+":"+proxy.Pass+"@"+proxy.Host+":"+proxy.Port
		cRes:=u.CheckProxy(proxyStr)
		
		if(!cRes){
			logs.Error(cRes)
			pxw.Replaceproxy(proxy.Host)
			//disable proxy
			u.DisableProxydb(proxy.Host,proxy.Port)
			continue;
		}

		_, err := u.Save(proxy)
		if err != nil {
			logs.Error(err)
			return err
		}
	}
	return nil
}

var CheckURL = "https://httpbin.org/get"
// check whether a string is work
//https://github.com/titanhw/go-proxy-checker/blob/master/core/checker.go
func (u *Proxy) CheckProxy(proxy string) bool {
	if strings.TrimSpace(proxy) == "" {
		return false
	}
	// no protocol, add //
	if !strings.Contains(proxy, "//") {
		proxy = "//" + proxy
	}
	
	// get resources from pool and release after operations
	// request := requestPool.Get().(*gorequest.SuperAgent)
	// resp := resultPool.Get().(map[string]interface{})
	var resp map[string]interface{}
	// defer requestPool.Put(request)
	// defer resultPool.Put(resp)
	// do the Request
	request := gorequest.New()
	// auth := username+":"+password
    // basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	// .Set("Proxy-Authorization", basicAuth)
	// logs.Info(proxy)
	_, _, errors := request.Proxy(proxy).Get(CheckURL).EndStruct(&resp)
	logs.Error(errors)
	return errors == nil 
	//return false
	//logs.Info(resp)
	
	//return strings.Contains(proxy, resp["origin"].(string))
}

//check proxy enable in google
func (u *Proxy) CheckGoogleProxy(proxy string, types string) bool {
	var CheckproxyURL = ""
	if strings.TrimSpace(proxy) == "" {
		return false
	}
	if(types=="google"){
		CheckproxyURL = "https://www.google.com"
	}else{
		CheckproxyURL = "https://www.bing.com"
	}
	// no protocol, add //
	if !strings.Contains(proxy, "//") {
		proxy = "//" + proxy
	}
	
	// get resources from pool and release after operations
	// request := requestPool.Get().(*gorequest.SuperAgent)
	// resp := resultPool.Get().(map[string]interface{})
	//var resp map[string]interface{}

	request := gorequest.New()

	resp, _, errors := request.Proxy(proxy).Get(CheckproxyURL).End()
	
	if errors != nil {
	logs.Error(errors)
		return false
	}
	logs.Info(resp.StatusCode)
	return resp.StatusCode==200
}
//update proxy
func (u *Proxy)Updateproxy()(error){
	pxw :=u.Getproxytype()
	return pxw.Updateproxy()
}
//get proxy from local database
func (u *Proxy)GetProxydb()([]Proxy,error){
	o := orm.NewOrm()
	var proxylist []Proxy
	_, err := o.QueryTable(u).Filter("available", 1).OrderBy("usetime", "-addtime").Limit(10).All(&proxylist)
	for _,proxy:=range proxylist{
		proxy.Usetime=time.Now()
		_, err := o.Update(&proxy)
		if err != nil {
			return proxylist,err
		}
	}
	return proxylist,err
}
//disable proxy from local database
func (u *Proxy)DisableProxydb(host string, port string)(int64,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// var proxyitem Proxy
	return qs.Filter("host", host).Filter("port", port).Update(orm.Params{
		"available": 0,
	})
}
//get proxy by id
func (u *Proxy)GetProxyById(id int64)(*Proxy,error){
	o := orm.NewOrm()
	var proxy Proxy
	err := o.QueryTable(new(Proxy)).Filter("id", id).One(&proxy)
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}
