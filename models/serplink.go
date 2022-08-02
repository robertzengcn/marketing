package models

import (
	// "fmt"
	"github.com/beego/beego/v2/client/orm"
)

type SerpLink struct {
	Id      int64      `json:"-" orm:"pk;auto"`
	Domain string `json:"domain" orm:"size(500)"`
	Link string `json:"link" orm:"size(1000)"`
	SearchrequestId *SearchRequest `json:"-" orm:"rel(fk);on_delete(do_nothing);column(searchrequest_id)"`
}

func (u *SerpLink) TableName() string {
	return "serplink"
}

// 设置引擎为 INNODB
func (u *SerpLink) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SerpLink))
	// create table
	// orm.RunSyncdb("default", false, true)
}
///save data to db
func (u *SerpLink) SavedataDb(serplink *SerpLink)(int64,error){
	o := orm.NewOrm()
	sid,oerr:=o.Insert(serplink)
	if(oerr!=nil){
		return 0,oerr
	}
	return sid,nil
}
///get serp link by request id
func(u *SerpLink)GetlistbyReqid(requestId int64)([]*SerpLink,int64,error){
	o := orm.NewOrm()
	var serpLinkVar SerpLink
	qs := o.QueryTable(&serpLinkVar)
	var serpLinkVarArr []*SerpLink
	num,err:=qs.Filter("searchrequest_id",requestId).All(&serpLinkVarArr)
	if(err!=nil){
		return nil,num,err
	}
	return serpLinkVarArr,0,nil
}