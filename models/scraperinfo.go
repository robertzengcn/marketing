package models
//save the info that scraper client get
import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
	"errors"
)

type ScraperInfo struct {
	Id int64 `orm:"pk;auto" json:"id"`
	Title string `orm:"size(1000)" json:"title"`
	Content string `orm:"type(text)" json:"content"`
	Url string `orm:"size(100)" json:"url"`
	Lang string `orm:"size(5)" json:"lang"`
	SocialTaskId *SocialTask `orm:"rel(fk);on_delete(do_nothing);column(socialtask_id)" json:"socialtask_id"`
	RecordTime  time.Time `orm:"auto_now;auto_now_add;type(datetime)" json:"record_time"`
	Usetime 	time.Time `orm:"null;type(datetime);column(usetime)" json:"use_time"`
	Processtime time.Time `orm:"null;type(datetime)" json:"process_time"`
}
func (u *ScraperInfo) TableName() string {
	return "scraperinfo"
}

// set table engineer INNODB
func (u *ScraperInfo) TableEngine() string {
	return "INNODB"
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(ScraperInfo))
	// create table
	// orm.RunSyncdb("default", false, true)
}

//save data to db
func (u *ScraperInfo) SavedataDb(scraperinfo *ScraperInfo) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ScraperInfo))
	var scraperinfoVar ScraperInfo
	// logs.Info("scraperinfo",scraperinfo)
	querys := qs.Filter("url", scraperinfo.Url)
	if(scraperinfo.SocialTaskId.Id>0){
		querys.Filter("socialtask_id",scraperinfo.SocialTaskId.Id)
	}
	err:=querys.One(&scraperinfoVar)
	logs.Info("scraperinfoVar",scraperinfoVar)
	logs.Info("socialtask_id",scraperinfo.SocialTaskId.Id)
	var sid int64
	var oerr error
	// logs.Error("err",err)
	if err == orm.ErrNoRows {
		sid, oerr = o.Insert(scraperinfo)

		if oerr != nil {
			return 0, oerr
		}
	}else{
		sid=scraperinfoVar.Id
	}
	return sid, nil
}
//get scraperinfo by taskId
func (u *ScraperInfo) GetScraperInfoByTaskId(taskId int64,limitNum int,uselimit bool,orderby string) ([]ScraperInfo, error) {
	if(orderby==""){
		orderby="usetime ASC"
	}
	o := orm.NewOrm()
	var scraperinfo []ScraperInfo
	oquery:= o.QueryTable(new(ScraperInfo)).Filter("socialtask_id", taskId).OrderBy("usetime").Limit(limitNum)
	if(uselimit){
		oquery.Filter("usetime",">=",time.Now().Add(-time.Hour*24).Format("2006-01-02 15:04:05"))
	}
	count, err :=oquery.All(&scraperinfo,"id","title","content","url","lang")
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	
	qs := o.QueryTable(u)
	currentTime := time.Now()
	//loop scraper info update use time
	for _, si := range scraperinfo {
		qs.Filter("id", si.Id).Update(orm.Params{
			"usetime": currentTime.Format("2006-01-02 15:04:05"),
		})
	}
	return scraperinfo, err
}
// find one scraper info by id
func (u *ScraperInfo)FindOne(id int64)(*ScraperInfo,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	scraperinfoVar:=ScraperInfo{}
	err :=qs.Filter("id",id).One(&scraperinfoVar)
	if(err!=nil){
		return nil,err
	}
	return &scraperinfoVar,nil
}
//update scraper info process time
func (u *ScraperInfo)UpdateProcesstime(id int64){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	currentTime := time.Now()
	qs.Filter("id", id).Update(orm.Params{
		"processtime": currentTime.Format("2006-01-02 15:04:05"),
	})
}

