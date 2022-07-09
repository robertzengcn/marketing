package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"marketing/utils"
	"os"
	"time"
	// "strconv"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"fmt"
)

var DefaultTask *Task

type Task struct {
	Id          int64       `orm:"pk;auto"`
	TaskName 	string      `orm:"size(150)" valid:"Required"`
	TaskStatus  *TaskStatus `orm:"rel(fk);on_delete(do_nothing)"`
	EmailTpl    *EmailTpl   `orm:"rel(fk);on_delete(do_nothing)" valid:"Required;"`
	CampaignId  *Campaign   `orm:"rel(fk);on_delete(do_nothing)" valid:"Required;"`
	CreatedTime time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdateTime  time.Time   `orm:"auto_now;type(datetime)"`
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
	orm.RegisterModelWithPrefix("mk_", new(Task))
	// create table
	// orm.RunSyncdb("default", false, true)
}
type Result struct {
	Runid 		int64
    Output          string
    Err error
}
///read search request json file and convert to json array
func (u *Task) Readfile(filename string) ([]SearchRequest, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var serequestarr []SearchRequest
	json.Unmarshal(byteValue, &serequestarr)
	return serequestarr, nil
}

///create task
func (u *Task) Createtask(task Task) (int64, error) {
	valid := validation.Validation{}
	b, verr := valid.Valid(&task)
    if verr != nil {
        // handle error
		return 0,verr
    }
    if !b {
        // validation does not pass
        var errMessage string
        for _, err := range valid.Errors {
			errMessage+=err.Key+":"+err.Message 
        }
		return 0,errors.New(errMessage)
    }
	o := orm.NewOrm()
	id, err := o.Insert(&task)
	if err != nil {
		return 0, err
	}
	return id,err	
}
///update task status
func (u *Task) Updatetaskstatus(taskId int64,taskStatusid int64)(error){
	o := orm.NewOrm()
	task := Task{Id: taskId}
	taskstatusModel:=TaskStatus{}
	taskStatusVar,statusErr:=taskstatusModel.GetOne(taskStatusid)
	if(statusErr!=nil){
		return errors.New("task status error")
	}
	terr:=o.Read(&task) 
	if(terr== nil) {
		task.TaskStatus=taskStatusVar
		if _, err := o.Update(&task); err != nil {
			return err //update failure
		}
	}else{
		return terr
	}
	return nil
}
///find one task by task id
func (u *Task)GetOne(taskId int64)(*Task,error){
	o := orm.NewOrm()
	task := Task{Id: taskId}
	err := o.Read(&task)
	if(err!=nil){
		return nil, err	
	}else{
		return &task,nil
	}
}
///start a task
func (u *Task)Starttask(taskId int64){
	taskrunModel:=TaskRun{}
	runid,runErr:=taskrunModel.CreateRun(taskId)
	if(runErr!=nil){
		logs.Error(runErr)
		return 
	}
	
	TaskdetailModel:=TaskDetail{}
	taskdetailVar,terr:=TaskdetailModel.Gettaskdetail(taskId)
	if(terr!=nil){		
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:terr})
		return
	}
	if(len(taskdetailVar.Taskkeyword)<=0){
		// return errors.New("keyword empty")
		
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:errors.New("keyword empty")})
	}

	gHost,gherr:=beego.AppConfig.String("googlescrape::host")
	if(gherr!=nil){
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:gherr})
		return
	}
	gPort,gperr:=beego.AppConfig.String("googlescrape::port")
	if(gperr!=nil){
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:gperr})
		return
	}
	gUser,gerr:=beego.AppConfig.String("googlescrape::user")
	if(gerr!=nil){
		// logs.Error(gerr)
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:gerr})
		return
	}
	gPass,gperr:=beego.AppConfig.String("googlescrape::pass")
	if(gperr!=nil){
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:gperr})
		return
	}
	conn, cerr := utils.Connect(gHost+":"+gPort, gUser, gPass)
	if cerr != nil {
		// logs.Error(cerr)
		u.Handletaskerror(&Result{Runid: runid,Output: "",Err:cerr})
		return
	}
	out := make(chan []byte)
    // errs := make(chan error)

	//create file over ssh
	// filename:=strconv.FormatInt(taskdetailVar.Id,10)
	keywordfile:="/app/GoogleScraper/"+taskdetailVar.TaskFilename+".txt"
	createfileCmd:="echo "+taskdetailVar.Taskkeyword+" > "+keywordfile
		
	// cmdArgs := []string{"-h"}
	logs.Info(createfileCmd)
	output, err := conn.SendCommands(createfileCmd)
	u.Handletaskerror(&Result{Runid: runid,Output: string(output),Err:err})
	if err != nil {
		logs.Error(err)	
		return
	}
	
	// out <-output
	// close(out)
	outputFile:="/app/GoogleScraper/"+taskdetailVar.TaskFilename+"-output.json"
	logs.Info(outputFile)
	keywordCom:="GoogleScraper -m selenium --sel-browser chrome --browser-mode headless --keyword-file "+keywordfile+" --num-workers 5 --output-filename "+outputFile+" -v debug"

	logs.Info(keywordCom)
	kout,kerr:=conn.SendCommands(keywordCom)
	logs.Info(string(kout))
	// out<-kout
	if(kerr!=nil){
		logs.Error(kerr)
		u.Handletaskerror(&Result{Runid: runid,Output: string(<-out),Err:kerr})
		return 
	}
	
}

///handle error during run task
func (u *Task)Handletaskerror(res *Result)(error){
	taskRunM:=TaskRun{}
	taskRunVar,rErr:=taskRunM.GetOne(res.Runid)	
	if(rErr!=nil){
		return rErr
	}
	// f, err := os.Create("/go_workspace/src/log/"+taskRunvar.Logid+".log.txt")
	// if(err!=nil){
	// 	return err
	// }
	// n3, err := f.WriteString(res.Output)
	// if(err!=nil){
	// 	return err
	// }
	logs.EnableFuncCallDepth(true)
	// logs.SetLogger(logs.AdapterMultiFile, `{"filename":"file.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	// logs.SetLogger(logs.AdapterFile,`{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	// logs.SetLogger(logs.AdapterMultiFile, `{"filename":"test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	fileName:=taskRunVar.Logid+".log"
	path:="/go_workspace/src"
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(
		`{"filename":"%s/log/%s", "daily":true,"rotate":true}`, path, fileName))
	if(res.Err!=nil){
		logs.Error(res.Err)
	}
	logs.Info(res.Output)
	
	// f.Sync()
	return nil
}


