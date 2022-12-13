package job
import (
	"github.com/beego/beego/v2/task"	
	"github.com/beego/beego/v2/core/logs"
	"marketing/models"
	"strconv"
	"context"
)
func InitTask(){

	cvtk1 := task.NewTask("tk1", "0 0 2 * * *", func(ctx context.Context) error { 
		
		logs.Info("run task per day")
		scheduleModel:=models.Schedule{}
		schVar, schErr:=scheduleModel.Findonebycyc("d")
		if(schErr!=nil){
			logs.Error(schErr)	
			return schErr
		}
		staId,staerr:=scheduleModel.Createtask(schVar.Id)
		if(staerr!=nil){
			logs.Error(staerr)	
			return staerr
		}
		logs.Info("task create with id"+strconv.FormatInt(staId, 10))
		TaskModel := models.Task{}

		go TaskModel.Starttask(staId)
		// TaskModel.Updatetaskstatus(staId, 3)
		return nil
		
	 })
	task.AddTask("tk1", cvtk1)

	// cvtk2 := task.NewTask("tk2", "0 0 1 * * *", func(ctx context.Context) error { 
	// 	logs.Info("start to get adult keywords")
	// 	KeywordModel:=models.Keyword{}
	// 	return KeywordModel.Getsexkeyword()
	// })
	// task.AddTask("tk2", cvtk2)
	//get sex toy keyword
	cvtk3 := task.NewTask("tk3", "0 30 1 * * *", func(ctx context.Context) error { 
		logs.Info("start to get sex toy keywords")
		KeywordModel:=models.Keyword{}
		_,kerr:=KeywordModel.Getkeywordapi()
		return kerr
	})
	task.AddTask("tk3", cvtk3)
}