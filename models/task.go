package models

import (
    "encoding/json"
    "io/ioutil"
    "os"
	"github.com/beego/beego/v2/client/orm"
)
var DefaultTask *Task
type Task struct {
	Id      int64     `orm:"pk;auto"`
	TaskStatus *TaskStatus `orm:"rel(fk);on_delete(do_nothing)"`
}

func (u *Task) TableName() string {
	return "task"
}

// set engineer as INNODB
func (u *Task) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Task))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///read search request json file and convert to json array
func (u *Task) Readfile(filename string)([]SearchRequest,error){
	jsonFile, err := os.Open(filename)
	if err != nil {
        return nil,err
    }	
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
    var serequestarr []SearchRequest
	json.Unmarshal(byteValue, &serequestarr)
	return 	serequestarr,nil
}




