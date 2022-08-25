package models
import (
	"github.com/beego/beego/v2/client/orm"
	"time" 
)

type Proxy struct {
	Id int64 `orm:"pk;auto"`
	Host string  `orm:"size(350)" valid:"Required"`
	Port string `orm:"size(6)" valid:"Required"`
	User string `orm:"size(350)"`
	Pass string `orm:"size(350)"`
	Protocol string `orm:"size(10)"`
	Available int `orm:"size(1);default(1);description(this mean status of the proxy)"`
	Addtime time.Time `orm:"auto_now_add;type(datetime)"`
	Checktime time.Time `orm:"null;type(datetime)"`	
}
///defined table name
func (u *Proxy) TableName() string {
	return "proxy"
}
func init() {
	orm.RegisterModelWithPrefix("mk_", new(Proxy))
}
// set engineer as INNODB
func (u *Proxy) TableEngine() string {
	return "INNODB"
}
