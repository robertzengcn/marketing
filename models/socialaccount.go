package models

import (
	"time"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)
var DefaultSocialAccount *SocialAccount
type SocialAccount struct {
	Id         int64  `orm:"pk;auto"`
	CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	UserName   string  `orm:"size(100)" json:"username"`
	PassWord   string  `orm:"size(100)" json:"password"`
	SocialplatformId *SocialPlatform `orm:"rel(fk);on_delete(do_nothing);column(socialplatform_id)" json:"socialplatform_id"`
	Stauts    int8   `orm:"default(1)"` // 1:active 2:inactive
	Createtime time.Time `orm:"auto_now;type(datetime)"`

}

func (u *SocialAccount) TableName() string {
	return "social_account"
}

// set engine as INNODB
func (u *SocialAccount) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SocialAccount))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///create social account data
func (u *SocialAccount)CreateSocialAccount(campaignId int64, userName string, passWord string, socialplatformId int64,accountname string,phoneNumber string,email string) (int64, error) {
	o := orm.NewOrm()
	socialAccount := SocialAccount{CampaignId: &Campaign{CampaignId: campaignId}, UserName: userName, PassWord: passWord, SocialplatformId: &SocialPlatform{Id: socialplatformId}}
	id, err := o.Insert(&socialAccount)
	if(err!=nil){
		//create social profile
		socialProModel:=new(SocialProfile)
		_,serr:=socialProModel.CreateSocialProfile(accountname,phoneNumber,email)
		if(serr!=nil){
			return id, serr
		}
	}
	return id, err
}



