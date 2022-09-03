package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type SearchRequest struct {
	Id      int64      `json:"-" orm:"pk;auto"`
	Query   string     `json:"query" orm:"size(100)"`
	Results *[]SerpLink `json:"results" orm:"-"`
	Taskrunid int64 `json:"-" orm:"column(task_run_id)"`
}

func (u *SearchRequest) TableName() string {
	return "searchrequest"
}

// set enginer as INNODB
func (u *SearchRequest) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SearchRequest))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///save search resuest list to db
func (u *SearchRequest) Savesrlist(secreq []SearchRequest,taskrunId int64) error {
	// o := orm.NewOrm()
		
	for _, item := range secreq {
		var searchreq SearchRequest
		searchreq.Query = item.Query
		searchreq.Results = item.Results
		searchreq.Taskrunid=taskrunId
		_, serr:=searchreq.SavedataDb(&searchreq)
		if(serr!=nil){
			fmt.Println(serr)
			break;
		}
		// sid, err := o.Insert(&searchreq)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// 	// break
		// }
		// for _, ritem := range *item.Results {
		// 	var serlink SerpLink
		// 	serlink.Domain=ritem.Domain
		// 	serlink.Link=ritem.Link
		// 	serlink.SearchrequestId=sid
		// 	serlink.SavedataDb(&serlink)
		// }
	}
	return nil
}

///save search request to database
func(u *SearchRequest)SavedataDb(searchreq *SearchRequest)(int64,error){
	o := orm.NewOrm()
	sid, err := o.Insert(searchreq)
	if(err!=nil){
		return 0,err
	}
	for _, ritem := range *searchreq.Results {
		var serlink SerpLink
		serlink.Domain=ritem.Domain
		serlink.Link=ritem.Link
		serlink.SearchrequestId=searchreq
		_,lerr:=serlink.SavedataDb(&serlink)
		if(lerr!=nil){
			return 0,lerr
		}
	}
	return sid,nil

}
///get request by run id
func(u *SearchRequest)Getrequestrunid(taskrunid int64)([]*SearchRequest,error){
	var srList []*SearchRequest
	o := orm.NewOrm()
	searchreq := SearchRequest{}
	_, err :=o.QueryTable(&searchreq).Filter("task_run_id", taskrunid).All(&srList)
	if(err!=nil){
		return nil, err
	}
	return srList,nil
}