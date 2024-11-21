package models

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	"marketing/utils"
)

type SocialtaskTagList struct {
	Id         int64 `orm:"pk;auto"`
	Tag        *Tag	`orm:"rel(fk);on_delete(do_nothing);column(mk_tag_id)"`
	SocialTask *SocialTask `orm:"rel(fk);on_delete(do_nothing);column(mk_social_task_id)"`
}

func (u *SocialtaskTagList) TableName() string {
	return "socialtask_tag_list"
}


func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(SocialtaskTagList))
	// create table
	// orm.RunSyncdb("default", false, true)
}

func (u *SocialtaskTagList) TableEngine() string {
	return "INNODB"
}



//get tag id by social task id
func (u *SocialtaskTagList) GetTagIdBySocialTaskId(socialTaskId int64) (tag []*Tag, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("mk_tag.tag_id,mk_tag.tag_name").From("mk_socialtask_tag_list").LeftJoin("mk_tag").On("mk_socialtask_tag_list.mk_tag_id = mk_tag.tag_id").Where("mk_socialtask_tag_list.mk_social_task_id = ?")
	sql := qb.String()

	// 执行 SQL 语句
	o := orm.NewOrm()
	_, qerr := o.Raw(sql, socialTaskId).QueryRows(&tag)
	if qerr != nil {
		logs.Error(qerr)
		return nil, qerr
	}
	// logs.Info(tag)
	return tag, nil
}

//update tag by task id
func (u *SocialtaskTagList) UpdateTagBySocialTaskId(socialTaskId int64, tagId int64) (id int64, err error) {
	o := orm.NewOrm()
	var socialtaskTagList SocialtaskTagList
	err = o.QueryTable(new(SocialtaskTagList)).Filter("social_task_id", socialTaskId).Filter("mk_tag_id", tagId).One(&socialtaskTagList)
	if err != nil {
		if err == orm.ErrNoRows {
			socialtaskTagList.SocialTask = &SocialTask{Id: socialTaskId}
			socialtaskTagList.Tag = &Tag{TagId: tagId}
		id,err:=o.Insert(&socialtaskTagList)
		return id,err
			}
		return 0, err	 
	}	
	// id, err = o.Update(&socialtaskTagList)
	return socialtaskTagList.Id, err
}
//update social task tags
func (u *SocialtaskTagList) UpdateSocialTaskTags(socialTaskId int64, tagIds []int64) (id int64, err error) {
	o := orm.NewOrm()
	//delete old tags
	// _, err = o.QueryTable(new(SocialtaskTagList)).Filter("mk_social_task_id", socialTaskId).Delete()
	// if err != nil {
	// 	return 0, err
	// }
	//get old tag ids by social task id
	// var oldTagIds []int64
	// o.QueryTable(new(SocialtaskTagList)).Filter("mk_social_task_id", socialTaskId).All(&oldTagIds)
	oldtagArr,_:=u.GetTagIdBySocialTaskId(socialTaskId)
	for _, oldTag := range oldtagArr {
		oexist := utils.ContainsType(tagIds, oldTag.TagId)
		if !oexist {
			u.DeleteSocialTaskTags(socialTaskId, oldTag.TagId)
		}
	}
	// logs.Info("tagids follow:")
	// logs.Info(tagIds)
	//insert new tags
	for _, tagId := range tagIds {
		logs.Info(tagId)
		var socialtaskTagList SocialtaskTagList
		socialtaskTagList.SocialTask = &SocialTask{Id: socialTaskId}
		socialtaskTagList.Tag = &Tag{TagId: tagId}
		// logs.Info(socialtaskTagList)
		id, err = o.Insert(&socialtaskTagList)
		if err != nil {
			return 0, err
		}
	}
	return id, err
}
//delete social task tags by social task id and tag id
func (u *SocialtaskTagList) DeleteSocialTaskTags(socialTaskId int64, tagId int64) (id int64, err error) {
	o := orm.NewOrm()
	//delete old tags
	_, err = o.QueryTable(new(SocialtaskTagList)).Filter("mk_social_task_id", socialTaskId).Filter("mk_tag_id", tagId).Delete()
	if err != nil {
		return 0, err
	}
	return id, err
} 

