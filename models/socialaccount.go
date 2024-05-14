package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"

	// "github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

var DefaultSocialAccount *SocialAccount

type SocialAccount struct {
	Id               int64           `orm:"pk;auto"`
	// CampaignId       *Campaign  `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	AccountId 	     *Account  `orm:"rel(fk);on_delete(do_nothing);column(account_id)"`
	UserName         string          `orm:"size(100)" json:"username"`
	PassWord         string          `orm:"size(100)" json:"password"`
	Socialplatform *SocialPlatform `orm:"rel(fk);on_delete(do_nothing);column(socialplatform_id)" json:"socialplatform_id"`
	Status           int8            `orm:"default(1)"` // 1:active 2:inactive
	// UseProxy         int8            `orm:"default(1)"` // whether to use proxy
	// Proxy            *SocialProxy    `orm:"rel(fk);on_delete(do_nothing);column(proxy_id)" json:"proxy_id"`
	Createtime       time.Time       `orm:"auto_now;type(datetime)"`
	Deletetime 		time.Time       `orm:"type(datetime);column(deletetime);null"` //set delete tag for data
}
type SocialAccountUpdate struct {
	AccountId int64
	UserName         string         
	PassWord         string          
	Socialplatform *SocialPlatform 
	Status           int8                      
	Proxy            *SocialProxy
	AccountName 		string   
	Phone			 string
	Email 			string
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
func (u *SocialAccount) CreateSocialAccount(soa SocialAccountUpdate) (int64, error) {

	//find social proxy by proxy id
	// sopModel := SocialProxy{}
	// var sop SocialProxy
	var err error
	// if soa.Proxy.Id != 0 {
	// 	sop, err = sopModel.GetSocialProxyById(soa.Proxy.Id)
	// 	// logs.Error(proxyId)
	// 	if err != nil {
	// 		return 0, errors.New("proxy not found")
	// 	}
	// }
	o := orm.NewOrm()
	socialAccount := SocialAccount{
		// CampaignId: &Campaign{CampaignId: campaignId},
		AccountId: 	  &Account{Id: soa.AccountId},
		UserName:         soa.UserName,
		PassWord:         soa.PassWord,
		Socialplatform: &SocialPlatform{Id: soa.Socialplatform.Id},
		// Proxy:            &sop,
		Status:           1,
	}
	//log.Info(socialAccount)
	id, err := o.Insert(&socialAccount)
	if err != nil {
		return id, err
	}
	//create social profile
	socialProModel := new(SocialProfile)
	_, serr := socialProModel.CreateSocialProfile(id,soa.AccountName, soa.Phone, soa.Email)
	if serr != nil {
		return id, serr
	}

	return id, err
}

//get social account relation with id
func (u *SocialAccount) GetSocialAccount(id int64,ownerId int64) (*SocialAccount, error) {
	o := orm.NewOrm()
	var socialAccount SocialAccount

	err := o.QueryTable(new(SocialAccount)).Filter("id", id).Filter("account_id", ownerId).Filter("deletetime__isnull", true).One(&socialAccount, "id", "user_name", "pass_word", "socialplatform_id", "status")
	if err != nil {
		return nil, err
	}
	return &socialAccount, nil
}
//list social account
func (u *SocialAccount) ListSocialaccount(ownerId int64,page int,size int,keywords string,socialplatformId int64,orderby string)([]*SocialAccount, error) {
	o := orm.NewOrm()
	var socialAccount []*SocialAccount
	qs := o.QueryTable(SocialAccount{}).Filter("deletetime__isnull", true).Filter("account_id", ownerId)
	
	if(len(keywords)>0){
		logs.Info(keywords)
		qs=qs.Filter("user_name__contains",keywords)
	}
	if(socialplatformId>0){
		qs=qs.Filter("socialplatform_id",socialplatformId)
	}
	
	if(len(orderby)>0){
		qs=qs.OrderBy(orderby)
	}
	
	// logs.Info(size)
	
	_,err:=qs.Limit(size, page).All(&socialAccount)
	if err != nil {
		return nil, err
	}
	//add platform name
	for _, v := range socialAccount {
		platformModel := SocialPlatform{}
		platform, err := platformModel.GetSocialPlatformById(v.Socialplatform.Id)
		if err != nil {
			return nil, err
		}
		v.Socialplatform = &platform
	
	}
	return socialAccount, nil

}
//update social account
func (u *SocialAccount) UpdateSocialAccount(id int64,ownerId int64,updateEntity SocialAccountUpdate) error {
	soEntity,_:=u.GetSocialAccount(id,ownerId)
	if(soEntity==nil){
		return errors.New("social account not found")
	}
	o := orm.NewOrm()
	socialAccount := SocialAccount{
		Id:               id,
		AccountId:        &Account{Id: ownerId},
		UserName:         updateEntity.UserName,
		PassWord:         updateEntity.PassWord,
		Socialplatform: &SocialPlatform{Id: updateEntity.Socialplatform.Id},
		// Proxy:            &SocialProxy{Id: updateEntity.Proxy.Id},
		Status : updateEntity.Status,
	}
	_, err := o.Update(&socialAccount)
	if err != nil {
		return err
	}
	spu:=SocialProfileUpdate{
		Name: updateEntity.AccountName,
		PhoneNumber: updateEntity.Phone,
		Email: updateEntity.Email,
	}	
	// logs.Info(spu)
	socialProModel := new(SocialProfile)
	_, serr := socialProModel.UpdateSocialProfile(id,spu)
	if(serr!=nil){
		return serr
	}
	return nil
}
//count social account number
func (u *SocialAccount) CountSocialaccount(ownerId int64,keywords string,socialplatformId int64) (int64, error) {
	o := orm.NewOrm()
	//select count(*) from social_accont where accountid=ownerId and (accountname like '%keywords%' )
	qs:= o.QueryTable(new(SocialAccount)).Filter("account_id", ownerId).Filter("deletetime__isnull", true)
	if(len(keywords)>0){
		qs=qs.Filter("user_name__contains",keywords)
	}
	if(socialplatformId>0){
		qs=qs.Filter("socialplatform_id",socialplatformId)
	}
	count,err:=qs.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}


//delete social account
func (u *SocialAccount) DeleteSocialAccount(id int64,ownerId int64) error {
	soEntity,_:=u.GetSocialAccount(id,ownerId)
	if(soEntity==nil){
		return errors.New("social account not found")
	}
	o := orm.NewOrm()
	_, err := o.QueryTable(new(SocialAccount)).Filter("id", id).Filter("account_id", ownerId).Update(orm.Params{
		"deletetime": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}