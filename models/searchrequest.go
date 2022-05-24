package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type SearchRequest struct {
	Id      int64      `json:"-" orm:"pk;auto"`
	Query   string     `json:"query" orm:"size(100)"`
	Results *[]SerpLink `json:"results" orm:"-"`
}

func (u *SearchRequest) TableName() string {
	return "searchrequest"
}

// 设置引擎为 INNODB
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
func (u *SearchRequest) Savesrlist(secreq []SearchRequest) error {
	// o := orm.NewOrm()
		
	for _, item := range secreq {
		var searchreq SearchRequest
		searchreq.Query = item.Query
		searchreq.Results = item.Results
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
