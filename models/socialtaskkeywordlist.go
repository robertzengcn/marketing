package models

import (
	// "github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	"marketing/utils"
)

type SocialtaskKeywordList struct {
	Id         int64 `orm:"pk;auto"`
	Keyword    *Keyword `orm:"rel(fk);on_delete(do_nothing);column(mk_keyword_id)"`
	SocialTask *SocialTask	`orm:"rel(fk);on_delete(do_nothing);column(mk_social_task_id)"`
}

func (u *SocialtaskKeywordList) TableName() string {
	return "socialtask_keyword_list"
}

func (u *SocialtaskKeywordList) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(SocialtaskKeywordList))
}

//get keyword id by social task id
func (u *SocialtaskKeywordList) GetKeywordIdBySocialTaskId(socialTaskId int64) (keyword []*Keyword, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("mk_keyword.id,mk_keyword.keyword").From("mk_socialtask_keyword_list").LeftJoin("mk_keyword").On("mk_socialtask_keyword_list.mk_keyword_id = mk_keyword.id").Where("mk_socialtask_keyword_list.mk_social_task_id = ?")
	sql := qb.String()
	// logs.Info(sql)
	// l.Println(sql)
	// 执行 SQL 语句
	o := orm.NewOrm()
	_, qerr := o.Raw(sql, socialTaskId).QueryRows(&keyword)
	if qerr != nil {
		return nil, qerr
	}

	return keyword, nil
}


//update keywords to social task
func (u *SocialtaskKeywordList) UpdateKeywordsToSocialTask(socialTaskId int64, keywordIds []int64) (err error) {

	//get old keywords
	oldKeyword, _ := u.GetKeywordIdBySocialTaskId(socialTaskId)
	if len(oldKeyword) > 0 {
		//loop old keyword, check if keywordIds has this item
		for _, kitem := range oldKeyword {
			oexist := utils.ContainsType(keywordIds, kitem.Id)
			if !oexist { //OLD not exist in new, delete it
				err = u.DeleteItemBySocialTaskIdAndKeywordId(socialTaskId, kitem.Id)
				if err != nil {
					return err
				}
			}
		}
	}
	o := orm.NewOrm()
	for _, keywordId := range keywordIds {
		socialtaskKeywordList := SocialtaskKeywordList{SocialTask: &SocialTask{Id: socialTaskId},
			Keyword: &Keyword{Id: keywordId}}

		chres, _ := u.CheckItemBySocialTaskIdAndKeywordId(socialTaskId, keywordId)
		if !chres {
			_, err = o.Insert(&socialtaskKeywordList)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

//delete item by social task id
func (u *SocialtaskKeywordList) DeleteItemBySocialTaskId(socialTaskId int64) (err error) {
	o := orm.NewOrm()
	// var socialtaskKeywordList SocialtaskKeywordList
	_, err = o.QueryTable(new(SocialtaskKeywordList)).Filter("mk_social_task_id", socialTaskId).Delete()
	if err != nil {
		return err
	}
	return nil
}

//delete item by social task id and keyword id
func (u *SocialtaskKeywordList) DeleteItemBySocialTaskIdAndKeywordId(socialTaskId int64, keywordId int64) (err error) {
	o := orm.NewOrm()
	// var socialtaskKeywordList SocialtaskKeywordList
	_, err = o.QueryTable(new(SocialtaskKeywordList)).Filter("mk_social_task_id", socialTaskId).Filter("mk_keyword_id", keywordId).Delete()
	if err != nil {
		return err
	}
	return nil
}

//check item by social task id and keyword id
func (u *SocialtaskKeywordList) CheckItemBySocialTaskIdAndKeywordId(socialTaskId int64, keywordId int64) (bool, error) {
	o := orm.NewOrm()
	// var socialtaskKeywordList SocialtaskKeywordList
	err := o.QueryTable(new(SocialtaskKeywordList)).Filter("mk_social_task_id", socialTaskId).Filter("mk_keyword_id", keywordId).One(&SocialtaskKeywordList{})
	if err != nil {
		if err == orm.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
