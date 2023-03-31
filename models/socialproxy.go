package models

import (
	// "strings"
	"time"
	"github.com/beego/beego/v2/client/orm"
)

type SocialProxy struct {
	Id int64 `orm:"pk;auto"`
	Url string `orm:"size(350)" valid:"Required"`
	Username string `orm:"size(350)"`
	Password string `orm:"size(350)"`
	Campaign *Campaign `orm:"rel(fk);column(campaign_id)" json:"campaign_id"`
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

