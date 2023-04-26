package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	//"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/beego/beego/v2/core/logs"
)

type Asocksproxy struct{}

type AsocksproxyResp struct {
	Success bool `json:"success"`
	Message struct {
		CountProxies int            `json:"countProxies"`
		Proxies      []Asockproxies `json:"proxies"`
	} `json:"message"`
}

type Asockproxies struct {
	Proxy       string `json:"proxy"`
	ExternalIP  string `json:"externalIp"`
	CountryCode string `json:"countryCode"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	CityName    string `json:"cityName"`
	Speed       int    `json:"speed"`
}

const Asocksurl = "https://api.asocks.com"

func (u *Asocksproxy) Gettoken() string {
	asockstoken := beego.AppConfig.DefaultString("asocks::token", "")
	return "Bearer " + asockstoken
}

//get proxy list from asocksproxy
func (u *Asocksproxy) Proxylist() ([]Proxy, error) {
	// asockstoken := beego.AppConfig.DefaultString("asocks::token", "")
	// var bearer = "Bearer " + asockstoken
	bearer := u.Gettoken()
	client := http.Client{}
	req, err := http.NewRequest("GET", Asocksurl+"/v2/proxy/ports", nil)
	req.Header.Set("Authorization", bearer)
	if err != nil {
		return nil, err
	}
	resp, rp := client.Do(req)
	if rp != nil {
		return nil, rp
	}
	defer resp.Body.Close()
	var result AsocksproxyResp
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//logs.Info(string(bodyBytes))
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	var proxylist []Proxy
	for _, proxy := range result.Message.Proxies {
		psa := strings.Split(proxy.Proxy, ":")
		if len(psa) < 2 {
			continue
		}
		proxylist = append(proxylist, Proxy{
			Host:        psa[0],
			Port:        psa[1],
			User:        proxy.Login,
			Pass:        proxy.Password,
			CountryCode: proxy.CountryCode,
			Protocol:    "http",
			Source:      "socks5",
			Available:   1,
		})
	}
	return proxylist, nil
}

//update proxy
func (u *Asocksproxy) Updateproxy() error {
	return u.CreateMultiproxy(20)
}

//create proxy at asocksproxy
func (u *Asocksproxy) Createproxy() error {
	bearer := u.Gettoken()
	// proxyNum := 10
	client := http.Client{}
	// body := []byte(`{
	// 	"name": "proxy created by api",
	// 	"auth_type_id": 2,
	// 	"country_code": "US"
	// 	"proxy_type_id":1,
	// 	"state":"",		
	// }`)
	postData := make(map[string]string)
    postData["name"] = "proxy created by api"
    postData["auth_type_id"] = "2"
    postData["country_code"] = "US"
	postData["proxy_type_id"]="1"
	postData["state"]=""
	postData["city"]=""
	postData["method_rotate"]=""
	postData["timeout"]=""
	postData["asn"]=""
	body := new(bytes.Buffer)
    w := multipart.NewWriter(body)
    for k,v :=  range postData{
        w.WriteField(k, v)
    }
    w.Close()
	//for i := 1; i < proxyNum; i++ {
		req, err := http.NewRequest("POST", Asocksurl+"/v2/proxy/ports", body)
		req.Header.Set("Authorization", bearer)
		req.Header.Set("Content-Type", w.FormDataContentType())
		if err != nil {
			return err
		}
		resp, rp := client.Do(req)
		if(resp.StatusCode!=200){
			return errors.New("response code is not 200,code is "+strconv.Itoa(resp.StatusCode))
		}
		data, _ := ioutil.ReadAll(resp.Body)
		if rp != nil {
			return rp
		}
		
		logs.Info(string(data))
		defer resp.Body.Close()
	//}
	return nil
}
//create server proxies
func(u *Asocksproxy) CreateMultiproxy(proxyNum int)error{
for i := 1; i < proxyNum; i++ {
	cerr:=u.Createproxy()
if(cerr!=nil){
	return cerr
}
}
return nil
}
//replace proxy
func(u *Asocksproxy) Replaceproxy(string)(error){
	return nil
}
