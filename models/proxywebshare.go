package models

import (
	"encoding/json"
	"net/http"
	"strconv"

	//"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
	//"io/ioutil"
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
)

type ProxyWebshare struct{}

// var DefaultProxyWebshare *ProxyWebshare
type WebshareResponse struct {
	Count    int              `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []WebshareResult `json:"results"`
}
type WebshareResult struct {
	CityName         string `json:"city_name"`
	CountryCode      string `json:"country_code"`
	CreatedAt        string `json:"created_at"`
	ID               string `json:"id"`
	LastVerification string `json:"last_verification"`
	Password         string `json:"password"`
	Port             int64  `json:"port"`
	ProxyAddress     string `json:"proxy_address"`
	Username         string `json:"username"`
	Valid            bool   `json:"valid"`
}
type ReplaceProxyreq struct {
	ToReplace ToReplace `json:"to_replace"`
	ReplaceWith []ReplaceWith `json:"replace_with"`
	DryRun bool `json:"dry_run"`
}
type ToReplace struct{
	Type    string `json:"type"`
	IPRange string `json:"ip_range"`
}
type ReplaceWith struct{
	Type        string `json:"type"`
	CountryCode string `json:"country_code"`
}

type ReplaceProxyrep struct {
	ID        int    `json:"id"`
	Reason    string `json:"reason"`
	ToReplace struct {
		Type    string `json:"type"`
		IPRange string `json:"ip_range"`
	} `json:"to_replace"`
	ReplaceWith []struct {
		Type        string `json:"type"`
		CountryCode string `json:"country_code"`
	} `json:"replace_with"`
	DryRun         bool        `json:"dry_run"`
	State          string      `json:"state"`
	ProxiesRemoved interface{} `json:"proxies_removed"`
	ProxiesAdded   interface{} `json:"proxies_added"`
	CreatedAt      string      `json:"created_at"`
	CompletedAt    interface{} `json:"completed_at"`
}

const WEBSHAREURL string = "https://proxy.webshare.io"

//get proxy list
func (u *ProxyWebshare) Proxylist() ([]Proxy, error) {
	websharetoken := beego.AppConfig.DefaultString("webshare::token", "")
	if(len(websharetoken)<1){
		return nil, errors.New("webshare token is empty")
	}
	//logs.Info(websharetoken)
	//send http get request
	client := http.Client{}
	req, err := http.NewRequest("GET", WEBSHAREURL+"/api/v2/proxy/list/?mode=direct", nil)
	req.Header.Set("Authorization", websharetoken)
	if err != nil {
		return nil, err
	}
	resp, rp := client.Do(req)
	if rp != nil {
		return nil, rp
	}
	defer resp.Body.Close()
	//decode response
	var result WebshareResponse
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//logs.Info(string(bodyBytes))
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var returnRes []Proxy
	i := 1
	for {
		req, err := http.NewRequest("GET", WEBSHAREURL+"/api/v2/proxy/list/?mode=direct&page="+strconv.Itoa(i), nil)
		req.Header.Set("Authorization", websharetoken)
		if err != nil {
			return nil, err
		}
		resp, rp := client.Do(req)
		if rp != nil {
			return nil, rp
		}
		defer resp.Body.Close()
		//decode response
		var result WebshareResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		for _, proxy := range result.Results {
			if !proxy.Valid {
				continue
			}
			returnRes = append(returnRes, Proxy{
				Host:        proxy.ProxyAddress,
				Port:        strconv.FormatInt(proxy.Port, 10),
				User:        proxy.Username,
				Pass:        proxy.Password,
				CountryCode: proxy.CountryCode,
				Protocol:    "http",
				Source:      "webshare",
				Available:   1,
			})
		}
		i++
		_,uerr:=url.ParseRequestURI(result.Next)
		if(len(result.Next) == 0||uerr!=nil){
			break
		}
		//logs.Info(returnRes)
		
	}
	return returnRes, nil
}

//create proxy
func (u *ProxyWebshare) Createproxy() error {
	return nil
}
func (u *ProxyWebshare)Updateproxy() (error){
	return nil
}
//replace proxy
func (u *ProxyWebshare)Replaceproxy(url string) (error){
	torep:=ToReplace{Type: "ip_range",IPRange: url}
	rpw:=ReplaceWith{Type: "country",CountryCode: "US"}
	
	reitem:=[]ReplaceWith{}
	reitem = append(reitem, rpw)

	rpreq:=ReplaceProxyreq{
		ToReplace:torep,
		ReplaceWith:reitem,
		DryRun:false,
	}
	Webshareurl:=WEBSHAREURL+"/api/v2/proxy/replace/"
	
	rpreqJson, err := json.Marshal(rpreq)
	if err != nil {
		return err
	}
	logs.Info(string(rpreqJson))
	req, errs := http.NewRequest("POST", Webshareurl, bytes.NewBuffer(rpreqJson))	
	websharetoken := beego.AppConfig.DefaultString("webshare::token", "")
	req.Header.Add("Authorization", websharetoken)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Custom-Header", "myvalue")
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
     
		return err
    }
	
    defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if errs != nil {
		return errs
	}
	// logs.Info(string(body))
	if(resp.StatusCode!=201){
		return errors.New("replace proxy failure,response code is not 201,code is "+strconv.Itoa(resp.StatusCode))
	}
	cRep:=ReplaceProxyrep{}
	if jErr := json.Unmarshal(body, &cRep); jErr != nil {
		return jErr
	}
	return nil
}
