package models

import (
	"marketing/utils"

	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
)

type SocialAccountProxyList struct {
	Id              int64          `orm:"pk;auto"`
	SocialAccountId *SocialAccount `orm:"rel(fk);on_delete(do_nothing);column(social_account_id)"`
	ProxyId         *Proxy         `orm:"rel(fk);on_delete(do_nothing);column(proxy_id)"`
}

func (u *SocialAccountProxyList) TableName() string {
	return "social_account_proxy_list"
}

func (u *SocialAccountProxyList) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(SocialAccountProxyList))
}

//get proxy id by social account id
func (u *SocialAccountProxyList) GetProxyBySocialAccountId(socialAccountId int64) (proxy []*Proxy, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("mk_proxy.id,mk_proxy.host,mk_proxy.port,mk_proxy.user,mk_proxy.pass,mk_proxy.protocol").From("mk_social_account_proxy_list").LeftJoin("mk_proxy").On("mk_social_account_proxy_list.proxy_id = mk_proxy.id").Where("mk_social_account_proxy_list.social_account_id = ?")
	sql := qb.String()
	// logs.Info("sql here")
	// logs.Info(sql)
	o := orm.NewOrm()
	_, qerr := o.Raw(sql, socialAccountId).QueryRows(&proxy)
	if qerr != nil {
		return nil, qerr
	}
	return proxy, nil
}

//save social account id and proxy id into table
func (u *SocialAccountProxyList) SaveSocialAccountProxyList(socialAccountId int64, proxyId int64) (int64, error) {
	o := orm.NewOrm()
	u.SocialAccountId = &SocialAccount{Id: socialAccountId}
	u.ProxyId = &Proxy{Id: proxyId}
	id, err := o.Insert(u)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//update proxy to social account
func (u *SocialAccountProxyList) UpdateProxysToSocialAccount(socialAccountId int64, proxyIds []int64) (err error) {

	//get old keywords
	oldProxys, _ := u.GetProxyBySocialAccountId(socialAccountId)
	if len(oldProxys) > 0 {
		//loop old keyword, check if keywordIds has this item
		for _, kitem := range oldProxys {
			oexist := utils.ContainsType(proxyIds, kitem.Id)
			if !oexist { //OLD not exist in new, delete it
				err = u.DeleteItemByAccountIdPrxoxyId(socialAccountId, kitem.Id)
				if err != nil {
					return err
				}
			}
		}
	}
	o := orm.NewOrm()
	for _, proId := range proxyIds {
		

		chres, _ := u.CheckItemBySocialAccountIdAndProxyId(socialAccountId, proId)
		if !chres {
			socialproxyListItem := SocialAccountProxyList{SocialAccountId: &SocialAccount{Id: socialAccountId},
			ProxyId: &Proxy{Id: proId}}
			_, err = o.Insert(&socialproxyListItem)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

//delete proxy id by social account id and proxy id
func (u *SocialAccountProxyList) DeleteItemByAccountIdPrxoxyId(socialTaskId int64, keywordId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(new(SocialAccountProxyList)).Filter("social_account_id", socialTaskId).Filter("proxy_id", keywordId).Delete()
	if err != nil {
		return err
	}
	return nil
}
//check item by social account id and proxy id
func (u *SocialAccountProxyList) CheckItemBySocialAccountIdAndProxyId(socialTaskId int64, proxyId int64) (bool, error) {
	o := orm.NewOrm()
	exist := o.QueryTable(new(SocialAccountProxyList)).Filter("social_account_id", socialTaskId).Filter("proxy_id", proxyId).Exist()
	return exist, nil
}
