package job

import (
	"context"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/task"
	"marketing/models"
	"strconv"
)

func InitTask() {

	cvtk1 := task.NewTask("tk1", "0 0 */2 * * *", func(ctx context.Context) error {

		logs.Info("run task every two hour per day")
		scheduleModel := models.Schedule{}
		schVar, schErr := scheduleModel.Findallschedule("d")
		if schErr != nil {
			logs.Error(schErr)
			return schErr
		}
		for _, v := range schVar {
			staId, staerr := scheduleModel.Createtask(v.Id,1)
			if staerr != nil {
				logs.Error(staerr)
				return staerr
			}
			logs.Info("task create with id" + strconv.FormatInt(staId, 10))
			TaskModel := models.Task{}

			TaskModel.Starttask(staId, "google")
			TaskModel.Starttask(staId, "bing")
		}
		// TaskModel.Updatetaskstatus(staId, 3)
		return nil

	})
	task.AddTask("tk1", cvtk1)


	cvtk2 := task.NewTask("tk2", "0 0 4 * * *", func(ctx context.Context) error {

		logs.Info("run task at 4 per day")
		scheduleModel := models.Schedule{}
		schVar, schErr := scheduleModel.Findonebycyc("e")
		if schErr != nil {
			logs.Error(schErr)
			return schErr
		}
		staId, staerr := scheduleModel.Createtask(schVar.Id,1)
		if staerr != nil {
			logs.Error(staerr)
			return staerr
		}
		logs.Info("task create with id" + strconv.FormatInt(staId, 10))
		TaskModel := models.Task{}

		TaskModel.Starttask(staId,"google")
		TaskModel.Starttask(staId,"bing")
		
		return nil

	})
	task.AddTask("tk2", cvtk2)
	//update prodxy on 1:00 everyday
	cvtk3 := task.NewTask("tk3", "0 0 1 * * *", func(ctx context.Context) error {
		proxymodel:=models.Proxy{}
		// perr:=proxymodel.Updateproxy()
		
		// if(perr!=nil){
		// 	logs.Error(perr)
		// }
		herr:=proxymodel.Handleproxy()
		if(herr!=nil){
			logs.Error(herr)
		}
		
		return nil
	})
	task.AddTask("tk3", cvtk3)


}
