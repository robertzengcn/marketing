package models

import (
	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
)

//defined account role list struct
type AccountRolesList struct {
	Id      int64     `orm:"pk;auto"`
	Account *Account  `orm:"rel(fk)"`
	Role    *AccountRole `orm:"rel(fk)"`
}


func (u *AccountRolesList) TableName() string {
	return "account_roles_list"
}

// 设置引擎为 INNODB
func (u *AccountRolesList) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	// orm.RegisterModelWithPrefix("mk_",new(AccountRolesList))
	// create table
	// orm.RunSyncdb("default", false, true)
}
//get account role by account id
func (u *AccountRolesList)GetAccountRoleByAccountId(accountId int64) (accountRole []*AccountRole, err error) {
	// l := logs.GetLogger()
	// var accountRole []*AccountRole
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("mk_account_role.name").From("mk_account_roles_list").LeftJoin("mk_account_role").On("mk_account_roles_list.role_id = mk_account_role.id").Where("mk_account_roles_list.account_id = ?")
	sql := qb.String()
	//l.Println(sql)
	// 执行 SQL 语句
	o := orm.NewOrm()
	_,qerr:=o.Raw(sql, accountId).QueryRows(&accountRole)
	if qerr != nil {
		return nil, qerr
	}
	return accountRole, nil
}