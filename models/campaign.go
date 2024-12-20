package models

import (
	// "fmt"
	"errors"
	// "fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	// "time"
)

var DefaultCampaign *Campaign

type Campaign struct {
	CampaignId   int64  `orm:"pk;auto"`
	CampaignName string `orm:"size(100)"`
	// EmailTpl   *EmailTpl	`orm:"rel(fk);on_delete(do_nothing)"`
	CampaignDescription string        `orm:"type(text);null"`
	Tags                []*Tag        `orm:"rel(m2m);rel_table(mk_campaign_tag_list)"`
	Types               *CampaignType `orm:"null;rel(fk);column(types_id)"` //the type of campaign, email, social
	Disable             int           `orm:"default(0)"`   //0: disabled, 1: enabled
	AccountId           *Account      `orm:"rel(fk);on_delete(do_nothing);column(account_id)"`
}

func (u *Campaign) TableName() string {
	return "campaign"
}

// 设置引擎为 INNODB
func (u *Campaign) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(Campaign))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///create campaign with name
func (u *Campaign) CreateCampaign(username string, tags string, types int16) (id int64, err error) {
	//check campaign type is exist
	var camType CampaignType
	o := orm.NewOrm()
	o.QueryTable(new(CampaignType)).Filter("campaign_type_id", types).One(&camType)
	if camType.CampaignTypeId <= 0 {
		return 0, errors.New("campaign type not exist")
	}
	var us Campaign
	us.CampaignName = username
	// us.Tags=tags
	us.Types = &camType
	id, err = o.Insert(&us)

	//add tags

	return id, err
}

/// show all campaign
func (u *Campaign) ListCampaign(start int, limitNum int, accountId int64) ([]Campaign, error) {
	o := orm.NewOrm()
	var cam []Campaign
	count, e := o.QueryTable(new(Campaign)).Filter("account_id", accountId).Limit(limitNum, start).All(&cam, "campaign_id", "campaign_name", "account_id", "types_id")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	campaignType := CampaignType{}
	//get campaigin Types
	for i, element := range cam {
		// logs.Info(element)
		camType, err := campaignType.GetCampaignTypeById(element.Types.CampaignTypeId)
		if err == nil {
			cam[i].Types = &camType
		}
	}

	//get Tags
	camTagModel := CampaignTagList{}
	for _, element := range cam {
	tags,terr:=camTagModel.GetCampaignTagByCampaignId(element.CampaignId)
	if(terr==nil){
		element.Tags=tags
	}
	}
	return cam, nil
}

/// find campaign by campaign id
func (u *Campaign) FindCambyid(id int64) (*Campaign, error) {
	o := orm.NewOrm()
	campaign := Campaign{CampaignId: id}
	err := o.Read(&campaign)

	return &campaign, err

}

/// list campaign by type
func (u *Campaign) ListCambytype(types int32, start int, limitNum int) ([]Campaign, error) {
	o := orm.NewOrm()
	var cam []Campaign
	count, e := o.QueryTable(new(Campaign)).Filter("types_id", types).Filter("disable", 0).Limit(limitNum, start).All(&cam, "campaign_id", "campaign_name", "types_id")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	return cam, nil
}

///count campagin number
func (u *Campaign) CountCampaign() (int64, error) {
	o := orm.NewOrm()
	count, e := o.QueryTable(new(Campaign)).Count()
	if e != nil {
		return 0, e
	}
	return count, nil
}
