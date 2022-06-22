package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/client/orm"
	"os/exec"
	// "fmt"
	"github.com/beego/beego/v2/core/logs"
	"bufio"
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
func (u *Task)Starttask(taskId int64)(error){
	cmdName := "ls"
	cmdArgs := []string{"-al"}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		logs.Error("Error creating StdoutPipe for Cmd")
		logs.Error(err)
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			logs.Info("start to run task get out | %s\n", scanner.Text())
			// fmt.Printf("docker build out | %s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		// os.Exit(1)
		logs.Error("Error starting Cmd")
		logs.Error(err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		// os.Exit(1)
		logs.Error("Error waiting for Cmd")
		logs.Error(err)
		return err
	}
	return nil
}

