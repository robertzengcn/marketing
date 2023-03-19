package models

import (
	"encoding/json"
	"net/http"
	"strconv"
	//"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	//"io/ioutil"
)

type ProxyWebshare struct{}
type WebshareResponse struct {
	Count    int              `json:"count"`
	Next     interface{}      `json:"next"`
	Previous interface{}      `json:"previous"`
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

const WEBSHAREURL string = "https://proxy.webshare.io"

//get proxy list
func (u *ProxyWebshare) Proxylist() ([]Proxy, error) {
	websharetoken := beego.AppConfig.DefaultString("webshare::token", "")
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
	page := 25
	for i := 0; i < result.Count; i += page {
		req, err := http.NewRequest("GET", WEBSHAREURL+"/api/v2/proxy/list/?mode=direct&limit="+strconv.Itoa(page)+"&offset="+strconv.Itoa(i), nil)
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
		//logs.Info(returnRes)
		return returnRes, nil
	}
	return nil, nil
}

//create proxy
func (u *ProxyWebshare) Createproxy() error {
	return nil
}
func (u *ProxyWebshare)Updateproxy() (error){
	return nil
}
