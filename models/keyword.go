package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
	"marketing/utils"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"github.com/beego/beego/v2/core/logs"
	"runtime"
	"path/filepath"
	"encoding/json"
)
type Keyword struct {
	Id int64  `orm:"pk;auto" json:"-"`
	Keyword string `orm:"size(150)" json:"keyword"`
	Tag string `orm:"size(200)" json:"tag"`
	// CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Created time.Time `orm:"null;auto_now_add;type(datetime)"`
	UsedTime time.Time `orm:"null;"`

}

const(ADULTSITE="adult_site")

func (u *Keyword) TableName() string {
	return "keyword"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(Keyword))
}

///get adult keyword from script
func (u *Keyword)Getsexkeyword()(error){
	gstatus, gserr := beego.AppConfig.Int("scrapesexkeyword::status")
	if gserr != nil {	
		return gserr
	}
	if(gstatus!=1){
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
	
	savefilemd:=utils.Md5V2(strconv.FormatInt(nsec, 10))
	savefile:=savefilemd+".json"
	remotejsonFile:="/app/workspace/"+savefile
	fetCommand:="scrapykeyword -o "+remotejsonFile
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
	localFilepath := apppath + "/output/"+savefile
	derr := conn.Downloadfile(sftpClient, remotejsonFile, localFilepath)
	if derr != nil {
		logs.Error(derr)
		// u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: derr})
		return derr
	}
	
	keyArr,ferr:=u.Readkeywordbyfile(localFilepath)
	
	if(ferr!=nil){
		return ferr
	}
	for _,x:= range keyArr {
		sid,sErr:=u.Savekeyworddb(x,ADULTSITE)
		if(sErr!=nil){
			logs.Error(sErr)
		}
		logs.Info(sid)
	}

	return nil
}
///read keyword list from json file
func (u *Keyword)Readkeywordbyfile(localFilepath string)([]Keyword,error){
	byteValue, ferr :=utils.ReadFile(localFilepath)
	if(ferr!=nil){
		return nil,ferr
	}
	var fetkeywordarr []Keyword
	json.Unmarshal(byteValue, &fetkeywordarr)
	return fetkeywordarr,nil
}
///save keyword to local db
func (u *Keyword)Savekeyworddb(keywordVar Keyword,tag string)(int64,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	keywordtar:=Keyword{
		Keyword: keywordVar.Keyword,
		Tag: tag,
	}
	err :=qs.Filter("keyword",keywordVar.Keyword).Filter("tag",tag).One(&keywordtar)
	
	if(err == orm.ErrNoRows){
		id, aerr := o.Insert(&keywordtar)
		if(aerr!=nil){
			return 0,aerr
		}
		return id, nil
	}

	return 0,nil
}
///find keyword by tags
func (u *Keyword)Getkeywordbytag(tagsArr []string,num int)([]*Keyword,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var keywordArrs []*Keyword
	_, kerr :=qs.Filter("tag__in", tagsArr).OrderBy("used_time").Limit(num).All(&keywordArrs)
	if(kerr!=nil){
		return nil,kerr
	}
	 currentTime := time.Now()
	
	for _,v:=range keywordArrs{
		logs.Info(v)
	qs.Filter("keyword", v.Keyword).Update(orm.Params{
		"used_time": currentTime.Format("2006.01.02 15:04:05"),
	})
	}

	return keywordArrs,nil
}


