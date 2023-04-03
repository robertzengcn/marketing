package models

import (
	// "strings"
	"time"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/core/logs"
)

type SocialProxy struct {
	Id int64 `orm:"pk;auto"`
	Url string `orm:"size(350)" valid:"Required"`
	Username string `orm:"size(350)" valid:"Required"`
	Password string `orm:"size(350)" valid:"Required"`
	Campaign *Campaign `orm:"rel(fk);column(campaign_id)" json:"campaign_id" valid:"Required"`
	Createtime time.Time `orm:"auto_now;type(datetime)"`
}

///defined table name
func (u *SocialProxy) TableName() string {
	return "social_proxy"
}
func init() {
	orm.RegisterModelWithPrefix("mk_", new(SocialProxy))
}
//save social proxy to database
func (u *SocialProxy) Save(proxy SocialProxy) (int64, error) {
	//valid data
	valid := validation.Validation{}
	b, berr := valid.Valid(&proxy)
    if berr != nil {
		logs.Error(berr)
        return 0,berr
    }
	if !b {
        // validation does not pass
        // blabla...
        for _, verr := range valid.Errors {
			logs.Error(verr.Key, verr.Message)
           return 0, verr
        }
    }
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var proxyitem SocialProxy
	err := qs.Filter("url", proxy.Url).Filter("username", proxy.Username).Filter("password", proxy.Password).Filter("campaign_id", proxy.Campaign.CampaignId).One(&proxyitem)
	// logs.Error(err)
	if err == orm.ErrNoRows {
		id, err := o.Insert(&proxy)
		return id, err
	}
	return 0, err
}	
//get social proxy by campaign id
func (u *SocialProxy) GetSocialProxyByCampaignId(campaignid int64) (SocialProxy, error) {
	o := orm.NewOrm()
	var socialproxy SocialProxy
	err := o.QueryTable(u).Filter("campaign_id", campaignid).One(&socialproxy)
	return socialproxy, err
}

