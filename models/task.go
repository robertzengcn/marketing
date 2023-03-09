package models

import (
	"encoding/json"
	"errors"
	// "io/ioutil"
	"marketing/utils"
	// "os"
	"time"
	// "strconv"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"path/filepath"
	"runtime"
	// "math/rand"
)

var DefaultTask *Task

type Task struct {
	Id         int64       `orm:"pk;auto"`
	TaskName   string      `orm:"size(150)" valid:"Required"`
	TaskStatus *TaskStatus `orm:"rel(fk);on_delete(do_nothing)"`
	// EmailTpl    *EmailTpl   `orm:"rel(fk);on_delete(do_nothing)" valid:"Required;"`
	CampaignId  *Campaign `orm:"rel(fk);on_delete(do_nothing)" valid:"Required;"`
	CreatedTime time.Time `orm:"null;auto_now_add;type(datetime)"`
	UpdateTime  time.Time `orm:"null;auto_now;type(datetime)"`
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
	Runid  int64
	Output string
	Err    error
}

///read search request json file and convert to json array
func (u *Task) Readfile(filename string) ([]SearchRequest, error) {
	// jsonFile, err := os.Open(filename)
	// if err != nil {
	// 	return nil, err
	// }
	// defer jsonFile.Close()
	// // read our opened xmlFile as a byte array.
	// byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue, _ := utils.ReadFile(filename)
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
		return 0, verr
	}
	if !b {
		// validation does not pass
		var errMessage string
		for _, err := range valid.Errors {
			errMessage += err.Key + ":" + err.Message
		}
		return 0, errors.New(errMessage)
	}
	o := orm.NewOrm()
	id, err := o.Insert(&task)
	if err != nil {
		return 0, err
	}
	return id, err
}

///update task status
func (u *Task) Updatetaskstatus(taskId int64, taskStatusid int64) error {
	o := orm.NewOrm()
	task := Task{Id: taskId}
	taskstatusModel := TaskStatus{}
	taskStatusVar, statusErr := taskstatusModel.GetOne(taskStatusid)
	if statusErr != nil {
		return errors.New("task status error")
	}
	terr := o.Read(&task)
	if terr == nil {
		task.TaskStatus = taskStatusVar
		if _, err := o.Update(&task); err != nil {
			return err //update failure
		}
	} else {
		return terr
	}
	return nil
}

///find one task by task id
func (u *Task) GetOne(taskId int64) (*Task, error) {
	o := orm.NewOrm()
	task := Task{Id: taskId}
	err := o.Read(&task)
	if err != nil {
		return nil, err
	} else {
		return &task, nil
	}
}

///start a task
func (u *Task) Starttask(taskId int64,searchenginer string) {
	u.Updatetaskstatus(taskId, 3)
	defer u.Updatetaskstatus(taskId, 4)
	taskrunModel := TaskRun{}
	runid, runErr := taskrunModel.CreateRun(taskId)
	if runErr != nil {
		logs.Error(runErr)
		return
	}

	serequestarr, scerr := u.Searchgoogle(taskId, runid,searchenginer)
	if scerr != nil {
		logs.Error(scerr)
		u.Handletaskerror(&Result{Runid: runid, Output: "", Err: scerr})
		return
	}

	//logs.Info(serequestarr)
	saerr := u.SaveSearchreq(serequestarr, runid)
	if saerr != nil {
		logs.Error(saerr)
		u.Handletaskerror(&Result{Runid: runid, Output: "", Err: saerr})
		return
	}

	logs.Info("task end")
	// u.Sendemail(runid)
}

//save search requese to db and fetch email
func (u *Task) SaveSearchreq(serequestarr []SearchRequest, runid int64) error {
	searchreqModel := SearchRequest{}
	serr := searchreqModel.Savesrlist(serequestarr, runid)
	if serr != nil {
		logs.Error(serr)
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: serr})
		return serr
	}
	logs.Info("start fetch email")
	fetchModel := FetchEmail{}
	fErr := fetchModel.Fetchtaskemail(runid)
	if fErr != nil {
		logs.Error(fErr)
		return fErr
	}
	return nil
}

//search keywords on google
func (u *Task) Searchgoogle(taskId int64, runid int64, searchenginer string) ([]SearchRequest, error) {
	TaskdetailModel := TaskDetail{}
	taskdetailVar, terr := TaskdetailModel.Gettaskdetail(taskId)
	if terr != nil {
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: terr})
		return nil, terr
	}
	if len(taskdetailVar.Taskkeyword) <= 0 {
		// return errors.New("keyword empty")
		u.Handletaskerror(&Result{Runid: runid, Output: "", Err: errors.New("keyword empty")})
	}

	//get host connect
	conn, cerr := u.Gethostconnect()
	if cerr != nil {
		// logs.Error(cerr)
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: cerr})
		return nil, cerr
	}
	workNum := beego.AppConfig.DefaultString("googlescrape::worrkernum", "1")
	// out := make(chan []byte)

	keywordfile := "/app/GoogleScraper/" + taskdetailVar.TaskFilename + ".txt"
	
	createfileCmd := "echo $'" + taskdetailVar.Taskkeyword + "' > " + keywordfile
	
	// cmdArgs := []string{"-h"}
	logs.Info(createfileCmd)
	output, err := conn.SendCommands(createfileCmd)
	u.Handletaskerror(&Result{Runid: runid, Output: string(output), Err: err})
	if err != nil {
		//logs.Error(err)
		return nil, err
	}

	outputFilename := taskdetailVar.TaskFilename + "-output.json"
	outputFile := "/app/GoogleScraper/" + outputFilename
	logs.Info(outputFile)
	nunPage := "10"
	// workNum := "2"
	keywordCom := "GoogleScraper -m selenium --sel-browser chrome --browser-mode headless --keyword-file " + keywordfile + " --num-workers " + workNum + " --output-filename " + outputFile + " --num-pages-for-keyword " + nunPage + " -v debug"
	if(searchenginer=="bing"){
		keywordCom+=" --search-engines \"bing\""
	}
	//logs.Info(keywordCom)
	kout, kerr := conn.SendCommands(keywordCom)
	logs.Info(string(kout))
	// out<-kout
	if kerr != nil {
		logs.Error(kerr)
		//u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: kerr})
		return nil, kerr
	}
	//read ssh file
	sftpClient, sftperr := conn.Createsfptclient()
	if sftperr != nil {
		logs.Error(sftperr)
		//u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: sftperr})
		return nil, sftperr
	}
	defer sftpClient.Close()
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	outPath := apppath + "/output/"
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err := os.Mkdir(outPath, 0755)
		if err != nil {
			logs.Error(err)
			//u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: err})
			return nil, err
		}
	}
	localFilepath := apppath + "/output/" + outputFilename
	derr := conn.Downloadfile(sftpClient, outputFile, localFilepath)
	if derr != nil {
		logs.Error(derr)
		//u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: derr})
		return nil, derr
	}

	serequestarr, rerr := u.Readfile(localFilepath)
	if rerr != nil {
		logs.Error(rerr)
		u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: rerr})
		return nil, rerr
	}
	return serequestarr, nil
}



//return search host connection
func (u *Task) Gethostconnect() (*utils.Connection, error) {
	gHost, gherr := beego.AppConfig.String("googlescrape::host")
	if gherr != nil {
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: gherr})
		return nil, gherr
	}
	gPort, gperr := beego.AppConfig.String("googlescrape::port")
	if gperr != nil {
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: gperr})
		return nil, gperr
	}
	gUser, gerr := beego.AppConfig.String("googlescrape::user")
	if gerr != nil {
		// logs.Error(gerr)
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: gerr})
		return nil, gerr
	}
	gPass, gperr := beego.AppConfig.String("googlescrape::pass")
	if gperr != nil {
		//u.Handletaskerror(&Result{Runid: runid, Output: "", Err: gperr})
		return nil, gperr
	}
	return utils.Connect(gHost+":"+gPort, gUser, gPass)

}

// ///start a task
// func (u *Task) StartBingtask(taskId int64) {
// 	u.Updatetaskstatus(taskId, 3)
// 	defer u.Updatetaskstatus(taskId, 4)
// 	taskrunModel := TaskRun{}
// 	runid, runErr := taskrunModel.CreateRun(taskId)
// 	if runErr != nil {
// 		logs.Error(runErr)
// 		return
// 	}

// 	serequestarr, scerr := u.Searchgoogle(taskId, runid,"bing")
// 	if scerr != nil {
// 		logs.Error(scerr)
// 		u.Handletaskerror(&Result{Runid: runid, Output: "", Err: scerr})
// 		return
// 	}

// 	logs.Info(serequestarr)

// 	searchreqModel := SearchRequest{}
// 	serr := searchreqModel.Savesrlist(serequestarr, runid)
// 	if serr != nil {
// 		logs.Error(serr)
// 		u.Handletaskerror(&Result{Runid: runid, Output: "", Err: serr})
// 		return
// 	}
// 	logs.Info("start fetch email after bing search")
// 	fetchModel := FetchEmail{}
// 	fErr := fetchModel.Fetchtaskemail(runid)
// 	if fErr != nil {
// 		logs.Error(fErr)
// 	}
// 	logs.Info("task end")
// 	// u.Sendemail(runid)
// }

///handle error during run task
func (u *Task) Handletaskerror(res *Result) error {
	taskRunM := TaskRun{}
	taskRunVar, rErr := taskRunM.GetOne(res.Runid)
	if rErr != nil {
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
	fileName := taskRunVar.Logid + ".log"
	// path := "/go_workspace/src"
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(
		`{"filename":"%s/log/%s", "daily":true,"rotate":true}`, apppath, fileName))
	if res.Err != nil {
		logs.Error(res.Err)
	}
	logs.Info(res.Output)

	// f.Sync()
	return nil
}

///send email for taskrun id
func (u *Task) Sendemail(tastrunId int64) error {
	//get CampaignId
	// taskrunModel := TaskRun{}
	// taskrun, terr := taskrunModel.GetOne(tastrunId)
	// if terr != nil {
	// 	return terr
	// }
	// taskModel := Task{}
	// task, taerr := taskModel.GetOne(taskrun.Task.Id)
	// if taerr != nil {
	// 	return taerr
	// }

	//get all email by task run id
	fetModel := FetchEmail{}
	femailslice, fnum, ferr := fetModel.Fetchallemail(tastrunId)
	if ferr != nil {
		return ferr
	}
	if fnum == 0 {
		return errors.New("not find email with task run id")
	}
	mailModel := MailLog{}
	emailser := EmailService{}
	for _, v := range femailslice {
		// logNum, _ := mailModel.Getemailcam(v.Email, task.CampaignId.CampaignId)
		// if logNum > 0 { //mail already exist in log
		// 	continue
		// }
		mbool, mErr := mailModel.Checkemailsend(v.Email, tastrunId)
		if mErr != nil {
			logs.Error(mErr)
		}
		if mbool {
			continue
		}
		//getmail account for send email
		// seremail,sererr:=emailser.GetEsbycam(task.CampaignId.CampaignId)
		// if(sererr!=nil){
		// 	return sererr
		// }
		// //get random email tpl
		// rand.Seed(time.Now().Unix())
		// chooseEm:=emArr[rand.Intn(len(emArr))]
		// toMail:= make([]string, 3)
		// toMail[0]=v.Email

		// //replace email content
		// chooseEm,reErr:=emailtplModel.Replacevar(chooseEm,v)
		// if(reErr!=nil){
		// 	logs.Error(reErr)
		// 	continue
		// }

		// //send email
		// serErr:=emailser.Sendemailtsl(seremail,toMail,chooseEm.TplTitle,chooseEm.TplContent)
		// if(serErr!=nil){
		// 	return serErr
		// }
		// maillogModel:=MailLog{Campaign: task.CampaignId,
		// 	Subject:chooseEm.TplTitle,
		// 	Content: chooseEm.TplContent,
		// 	Receiver: toMail[0],
		// 	TaskrunId: taskrun,
		//  }
		//  maillogModel.Addmaillog(maillogModel)
		emailser.Sendemailtask(v, tastrunId)
	}
	return nil
}
