package models

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"marketing/utils"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Keyword struct {
	Id      int64  `orm:"pk;auto" json:"id"`
	Keyword string `orm:"size(150)" json:"keyword"`
	//Tag string `orm:"size(200)" json:"tag"`
	Tag *Tag `orm:"null;rel(fk);column(tag_id)" json:"tag_id"`
	// CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Account  *Account  `orm:"rel(fk);on_delete(do_nothing);column(account_id)" json:"account_id"`
	Created  time.Time `orm:"null;auto_now_add;type(datetime)"`
	UsedTime time.Time `orm:"null;"`
}
type Adultapiresp struct {
	Status bool     `json:"status"`
	Msg    string   `json:"msg"`
	Code   int      `json:"code"`
	Data   []string `json:"data"`
}

const (
	ADULTSITE = "adult_site"
)

func (u *Keyword) TableName() string {
	return "keyword"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(Keyword))
}

///get adult keyword from script
func (u *Keyword) Getsexkeyword(accountId int64) error {
	gstatus, gserr := beego.AppConfig.Int("scrapesexkeyword::status")
	if gserr != nil {
		return gserr
	}
	if gstatus != 1 {
		return nil
	}
	gHost, gherr := beego.AppConfig.String("scrapesexkeyword::host")
	if gherr != nil {
		return gherr
	}
	gPort, gperr := beego.AppConfig.String("scrapesexkeyword::port")
	if gperr != nil {
		return gperr
	}
	gUser, guerr := beego.AppConfig.String("scrapesexkeyword::user")
	if guerr != nil {
		return guerr
	}
	gPass, gpserr := beego.AppConfig.String("scrapesexkeyword::pass")
	if gpserr != nil {
		return gpserr
	}

	conn, cerr := utils.Connect(gHost+":"+gPort, gUser, gPass)
	if cerr != nil {
		// logs.Error(cerr)
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: cerr})
		return cerr
	}
	now := time.Now()
	nsec := now.UnixNano()

	savefilemd := utils.Md5V2(strconv.FormatInt(nsec, 10))
	savefile := savefilemd + ".json"
	remotejsonFile := "/app/workspace/" + savefile
	fetCommand := "scrapykeyword -o " + remotejsonFile
	logs.Info(fetCommand)
	kout, kerr := conn.SendCommands(fetCommand)
	logs.Info(string(kout))
	if kerr != nil {
		logs.Error(kerr)
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: kerr})
		return kerr
	}
	sftpClient, sftperr := conn.Createsfptclient()
	if sftperr != nil {
		logs.Error(sftperr)
		return sftperr
	}
	defer sftpClient.Close()
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	localFilepath := apppath + "/output/" + savefile
	derr := conn.Downloadfile(sftpClient, remotejsonFile, localFilepath)
	if derr != nil {
		logs.Error(derr)
		// u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: derr})
		return derr
	}

	keyArr, ferr := u.Readkeywordbyfile(localFilepath)

	if ferr != nil {
		return ferr
	}
	for _, x := range keyArr {
		sid, sErr := u.Savekeyworddb(x, 0, accountId)
		if sErr != nil {
			logs.Error(sErr)
		}
		logs.Info(sid)
	}

	return nil
}

///get keyword from api
func (u *Keyword) Getkeywordapi(accountId int64) ([]Keyword, error) {

	siteurl, siteerr := config.String("mainsite::url")
	if siteerr != nil {
		return nil, siteerr
	}
	siteacc, accerr := config.String("mainsite::user")
	if accerr != nil {
		return nil, accerr
	}
	sitepass, passerr := config.String("mainsite::pass")
	if passerr != nil {
		return nil, passerr
	}
	url := siteurl + "/keywords/list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+utils.BasicAuth(siteacc, sitepass))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// panic(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	adultapires := Adultapiresp{}
	// var emailarr []string
	if jErr := json.Unmarshal(body, &adultapires); jErr != nil {
		return nil, jErr
	}
	if !adultapires.Status {
		return nil, errors.New(adultapires.Msg)
	}
	var fetkeywordarr []Keyword
	for _, eval := range adultapires.Data {
		keywordModel := Keyword{Keyword: strings.TrimSpace(eval)}
		fetkeywordarr = append(fetkeywordarr, keywordModel)
		u.Savekeyworddb(keywordModel, 0, accountId)
	}
	logs.Info(fetkeywordarr)
	return fetkeywordarr, nil
}

///read keyword list from json file
func (u *Keyword) Readkeywordbyfile(localFilepath string) ([]Keyword, error) {
	byteValue, ferr := utils.ReadFile(localFilepath)
	if ferr != nil {
		return nil, ferr
	}
	var fetkeywordarr []Keyword
	json.Unmarshal(byteValue, &fetkeywordarr)
	return fetkeywordarr, nil
}

///save keyword to local db
func (u *Keyword) Savekeyworddb(keywordVar Keyword, tagId int64, accountId int64) (int64, error) {
	if tagId > 0 {
		//check tagId is exist
		TagModel := Tag{}
		_, terr := TagModel.GetTagByTagId(tagId)
		if terr != nil {
			return 0, terr
		}
	}
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	keywordtar := Keyword{
		Keyword: keywordVar.Keyword,
		Tag:     &Tag{TagId: tagId},
		Account: &Account{Id: accountId},
	}
	err := qs.Filter("keyword", keywordVar.Keyword).Filter("tag_id", tagId).Filter("account_id", accountId).One(&keywordtar)
	if err != nil {
		if err == orm.ErrNoRows {
			id, aerr := o.Insert(&keywordtar)
			if aerr != nil {
				return 0, aerr
			}
			return id, nil
		}
		return 0, err
	} else {
		return keywordtar.Id, nil
	}

	// return 0, nil
}

///find keyword by tags
///tagsArr the tag want to find
///num the number of keyword want to find
func (u *Keyword) Getkeywordbytag(tagsArr []string, num int, accountId int64) ([]*Keyword, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var keywordArrs []*Keyword
	_, kerr := qs.Filter("tag__in", tagsArr).Filter("account_id", accountId).OrderBy("used_time").Limit(num).All(&keywordArrs)
	if kerr != nil {
		return nil, kerr
	}
	currentTime := time.Now()

	for _, v := range keywordArrs {
		// logs.Info(v)
		qs.Filter("keyword", v.Keyword).Update(orm.Params{
			"used_time": currentTime.Format("2006-01-02 15:04:05"),
		})
	}

	return keywordArrs, nil
}

//create keyword list from csv data result
func (u *Keyword) CreateRescsv(filepath string, accountId int64) ([]Keyword, error) {
	data, err := utils.Csvfilehandle(filepath)
	if err != nil {
		return nil, err
	}
	var keywordArrs []Keyword
	TagModel := Tag{}
	for i, line := range data {
		if i > 0 { // omit header line
			var rec Keyword
			for j, field := range line {
				if j == 0 {
					rec.Keyword = field
				} else if j == 1 {
					tagVar, tagerr := TagModel.Checktag(field, accountId)
					if tagerr != nil {
						rec.Tag = tagVar
					} else {
						rec.Tag = tagVar
					}
				}
			}
			keywordArrs = append(keywordArrs, rec)
		}
	}
	return keywordArrs, nil
}

//save keyword list
func (u *Keyword) Savekeyword(keywordArrs []Keyword, tag int64, accountId int64) ([]Keyword, error) {
	var keywordArrsNew []Keyword
	for _, v := range keywordArrs {
		_, aerr := u.Savekeyworddb(v, tag, accountId)
		if aerr != nil {
			keywordArrsNew = append(keywordArrsNew, v)
		}
	}
	return keywordArrsNew, nil
}
