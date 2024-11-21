package models
import (
	"github.com/beego/beego/v2/client/orm"
)

type SocialTaskType struct {
	TypeId int64 `orm:"pk;auto" json:"id"` 
	TypeName string `orm:"size(150)" json:"name"`
}
func (u *SocialTaskType) TableName() string {
	return "socialtask_type"
}

// 设置引擎为 INNODB
func (u *SocialTaskType) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(SocialTaskType))
}

//list social task type
func (u *SocialTaskType) ListSocialTaskType() (socialTaskType []*SocialTaskType, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	_, err = qs.All(&socialTaskType)
	return socialTaskType, err
}
//add social task type
func (u *SocialTaskType) AddSocialTaskType(socialTaskType string) (id int64, err error) {
	o := orm.NewOrm()
	var socialTaskTypeEntity SocialTaskType
	socialTaskTypeEntity.TypeName=socialTaskType
	id, err = o.Insert(&socialTaskTypeEntity)
	return id,err
}
//find social task type by id
func (u *SocialTaskType) FindSocialTaskTypeById(id int64) (*SocialTaskType,error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	socialTaskType:=SocialTaskType{}
	err := qs.Filter("type_id", id).One(&socialTaskType)
	return &socialTaskType, err
}