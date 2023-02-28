package models
import(
	"time"
	"github.com/beego/beego/v2/client/orm"
)
var DefaultApiauth *Apiauth

type Apiauth struct {
	Id         int64  `orm:"pk;auto"`
	UserName   string `orm:"size(100)"`
	Password   string `orm:"size(100)"`
	Updated time.Time `orm:"null;auto_now;type(datetime)"`
}

// set engine as INNODB
func (u *Apiauth) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Apiauth))
}

//get api data id by username and password
func (u *Apiauth)GetApiAuth(username string,password string) (int64, error) {
	o := orm.NewOrm()
	var apiauth Apiauth
	err := o.QueryTable(new(Apiauth)).Filter("UserName", username).Filter("Password", password).One(&apiauth, "Id")
	if err != nil {
		return 0, err
	}
	return apiauth.Id, nil
}