package tests

import (
	"marketing/models"
	"testing"
	"fmt"
)

func TestReadfile(t *testing.T) {
	// var filename string
	filename:="./20220509-threaded-results.json"
	res,rerr:=models.DefaultTask.Readfile(filename)
	if(rerr!=nil){
		panic(rerr.Error())	
	}
	fmt.Printf("%+v\n", res)
	

}