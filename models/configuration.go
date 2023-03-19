package models

import (
	 //"errors"
	 "time"
	  "github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
	ConfigurationId int64  `orm:"pk;auto"`
	ConfigurationKey string `orm:"size(50)"`
	ConfigurationValue string `orm:"size(300)"`
	Created time.Time `orm:"null;auto_now_add;type(datetime)"`
	Updated time.Time `orm:"null;auto_now;type(datetime)"`
}

func (u *Configuration) TableName() string {
	return "configuration"
}

// set mysql enginer as INNODB
func (u *Configuration) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Configuration))
}

//save configuration to database
func (u *Configuration)SaveConfiguration(configurationKey string, configurationValue string) (int64, error) {
	o := orm.NewOrm()
	var configvar Configuration
	configvar.ConfigurationKey = configurationKey
	configvar.ConfigurationValue = configurationValue
	return o.Insert(configvar)
}








