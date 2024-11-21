package models

import (
	// "github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type SocialProfile struct {
	Id          int64          `orm:"pk;auto"`
	Name        string         `orm:"size(100)"`
	PhoneNumber string         `orm:"size(100)"`
	Email       string         `orm:"size(100)"`
	Account     *SocialAccount `orm:"rel(fk);on_delete(do_nothing);column(social_account_id);unique"`
}
type SocialProfileUpdate struct {
	Name        string
	PhoneNumber string
	Email       string
}

// set engine as INNODB
func (u *SocialProfile) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(SocialProfile))
}

//create social SocialProfile Data
func (u *SocialProfile) CreateSocialProfile(socialaccountId int64, name string, phoneNumber string, email string) (int64, error) {
	o := orm.NewOrm()
	socialProfile := SocialProfile{
		Account: &SocialAccount{Id: socialaccountId},
		Name:    name, PhoneNumber: phoneNumber, Email: email}
	id, err := o.Insert(&socialProfile)
	return id, err
}

//update social profile data by account id
func (u *SocialProfile) UpdateSocialProfile(accountId int64, profileData SocialProfileUpdate) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// logs.Info("accountId",accountId)
	// logs.Info("name",profileData.Name)
	// logs.Info("phone_number",profileData.PhoneNumber)
	// sp:=SocialProfile{
	// 	Id:
	// }
	// _, err := o.Update(&sp)
	// acnum,acerr:=qs.Filter("account_id", 7).Count()
	// logs.Info("acnum",acnum)
	// if(acerr!=nil){
	// 	logs.Error("acerr",acerr)
	// }
	sp := SocialProfile{}
	qs.Filter("social_account_id", accountId).One(&sp)
	if sp.Id > 0 {
		//update social profile
		_, err := qs.Filter("social_account_id", accountId).Update(orm.Params{
			"name":         profileData.Name,
			"phone_number": profileData.PhoneNumber,
			"email":        profileData.Email,
		})
		return sp.Id, err
	} else {
		//insert social profile
		socialProfile := SocialProfile{
			Account:     &SocialAccount{Id: accountId},
			Name:        profileData.Name,
			PhoneNumber: profileData.PhoneNumber,
			Email:       profileData.Email}
		id, err := o.Insert(&socialProfile)
		return id, err
	}
}

//get social profile by social account id
func (u *SocialProfile) GetSocialProfileByAccountId(socialaccountId int64) (SocialProfile, error) {
	o := orm.NewOrm()
	profile := SocialProfile{}
	qs := o.QueryTable(new(SocialProfile))
	err := qs.Filter("social_account_id", socialaccountId).One(&profile, "Name", "Email", "PhoneNumber")

	// profile := SocialProfile{Account: &SocialAccount{Id: socialaccountId}}

	// err := o.Read(&profile, "Name","Email","PhoneNumber","Account")

	return profile, err
}
