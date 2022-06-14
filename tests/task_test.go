package tests

import (
	"marketing/models"
	"testing"
	"fmt"
	"runtime"
	"path/filepath"
)

func TestReadfile(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	filename:=apppath+"/tests/20220509-threaded-results.json"
	res,rerr:=models.DefaultTask.Readfile(filename)
	if(rerr!=nil){
		panic(rerr.Error())	
	}
	fmt.Printf("%+v\n", res)
	

}