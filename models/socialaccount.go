package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

var DefaultSocialAccount *SocialAccount

type SocialAccount struct {
	Id               int64           `orm:"pk;auto"`
	CampaignId       *Campaign       `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	UserName         string          `orm:"size(100)" json:"username"`
	PassWord         string          `orm:"size(100)" json:"password"`
	SocialplatformId *SocialPlatform `orm:"rel(fk);on_delete(do_nothing);column(socialplatform_id)" json:"socialplatform_id"`
	Stauts           int8            `orm:"default(1)"` // 1:active 2:inactive
	Proxy            *SocialProxy    `orm:"rel(fk);on_delete(do_nothing);column(proxy_id)" json:"proxy_id"`
	Createtime       time.Time       `orm:"auto_now;type(datetime)"`
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
	orm.RegisterModelWithPrefix("mk_", new(SocialAccount))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///create social account data
func (u *SocialAccount) CreateSocialAccount(campaignId int64, userName string, passWord string, socialplatformId int64, accountname string, phoneNumber string, email string, proxyId int64) (int64, error) {
	//find social proxy by proxy id
	sopModel := SocialProxy{}
	sop, err := sopModel.GetSocialProxyById(proxyId)
	logs.Error(proxyId)
	if err != nil {
		return 0, errors.New("proxy not found")
	}

	o := orm.NewOrm()
	socialAccount := SocialAccount{CampaignId: &Campaign{CampaignId: campaignId},
		UserName:         userName,
		PassWord:         passWord,
		SocialplatformId: &SocialPlatform{Id: socialplatformId},
		Proxy:            &sop,
		Stauts: 1,
	}
	//log.Info(socialAccount)
	id, err := o.Insert(&socialAccount)
	if err != nil {
		return id, err
	}
	//create social profile
	socialProModel := new(SocialProfile)
	_, serr := socialProModel.CreateSocialProfile(accountname, phoneNumber, email)
	if serr != nil {
		return id, serr
	}

	return id, err
}

///get social account relation with campaign use CampaignId
func  (u *SocialAccount)GetSocialAccountcam(id int64)(*SocialAccount,error){
	o := orm.NewOrm()
	var socialAccount SocialAccount
	
	err := o.QueryTable(new(SocialAccount)).Filter("campaign_id", id).Filter("stauts", 1).One(&socialAccount, "id", "campaign_id", "user_name", "pass_word", "socialplatform_id", "proxy_id")
	if err != nil {
		return nil, err
	}
	return &socialAccount, nil
}