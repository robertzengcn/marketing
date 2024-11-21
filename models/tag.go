package models
import (
	"github.com/beego/beego/v2/client/orm"
	// "reflect"
	"github.com/beego/beego/v2/core/logs"
)

type Tag struct {
	TagId int64 `orm:"pk;auto"`
	TagName string `orm:"size(150)"`
	AccountId *Account `orm:"rel(fk);column(account_id)"`
}
func (u *Tag) TableName() string {
	return "tag"
}

// 设置引擎为 INNODB
func (u *Tag) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(Tag))
}

//add tags by string
func (u *Tag) AddTagsByString(tags string,accountId int64) (id int64, err error) {
	//check account exist
	var account Account
	accountVar,_:=account.GetAccountbyid(accountId)
	if(accountVar.Id<=0){
		//acount not exist
		return 0,err
	}
	
	o := orm.NewOrm()	
	//check tag exist
	
	tagvar,_:=u.Checktag(tags,accountId)
	if(tagvar.TagId<=0){
		logs.Info("tag not exist")
		var tag Tag
		tag.TagName=tags
		tag.AccountId=&Account{Id:accountId}

		id, err = o.Insert(&tag)
		return id,err
	}else{
		logs.Info("tag exist")
		return tagvar.TagId,nil
	}

	
}
//check tag exist by string
func (u *Tag) Checktag(tag string,accountId int64)(*Tag,error){
	logs.Info("tag:",tag)
	logs.Info("accountId:",accountId)
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	tagentity:=Tag{}
	err :=qs.Filter("tag_name",tag).Filter("account_id",accountId).One(&tagentity)
	logs.Info("tagentity:",tagentity)
	return &tagentity,err
}
//list tag by account id
func (u *Tag) ListTagByAccountId(accountId int64) (tags []*Tag, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	_, err = qs.Filter("account_id", accountId).All(&tags)
	return tags, err
}
//get tag by tag id
func (u *Tag) GetTagByTagId(tagId int64) (tag *Tag, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	tag = &Tag{}
	err = qs.Filter("tag_id", tagId).One(tag)
	return tag, err
}