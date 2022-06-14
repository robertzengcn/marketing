package tests

import (
	"marketing/models"
	"testing"
	"fmt"
	"runtime"
	"path/filepath"
	beego "github.com/beego/beego/v2/server/web"
	// . "github.com/smartystreets/goconvey/convey"
	// log "github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}
///
func TestSavereuest(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	filename:=apppath+"/tests/20220509-threaded-results.json"
	res,rerr:=models.DefaultTask.Readfile(filename)
	if(rerr!=nil){
		panic(rerr.Error())	
	}
	searchreqModel:=models.SearchRequest{}
	serr:=searchreqModel.Savesrlist(res)
	if(serr!=nil){
		fmt.Println(serr)
	}

}