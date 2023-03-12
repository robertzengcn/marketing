package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
	"github.com/parnurzeal/gorequest"
	"strings"
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
	pxw := ProxyWebshare{}
	proarr,perr:=u.GetProxylist(&pxw)
	if(perr!=nil){
		return nil
	}
	//range proxy list,save to database
	for _, proxy := range proarr {
		_, err := u.Save(proxy)
		if err != nil {
			return err
		}
	}
	return nil
}
var CheckURL = "https://httpbin.org/get"
// check whether a string is work
//https://github.com/titanhw/go-proxy-checker/blob/master/core/checker.go
func CheckProxy(proxy string) bool {
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
	_, _, errors := request.Proxy(proxy).Get(CheckURL).EndStruct(&resp)
	if errors != nil {
		return false
	}
	return strings.Contains(proxy, resp["origin"].(string))
}

