package utils

import (
    beego "github.com/beego/beego/v2/server/web"
	"regexp"
	"os/exec"
	"github.com/beego/beego/v2/core/logs"
	"bufio"
	"crypto/md5"
	"fmt"
)

func Init(){
    beego.AddFuncMap("ValidEmail", ValidEmail) 
	beego.AddFuncMap("Contains", Contains)
    //add new function in here
}
//valid email valid
func ValidEmail(email string) bool{
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
///run command
func Runcommand(cmdName string,cmdArgs ...string)error{
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
			logs.Info(scanner.Text())
			// fmt.Printf("docker build out | %s\n", scanner.Text())
		}
	}()
	cmderrReader,cerr:=cmd.StderrPipe()
	if(cerr!=nil){
		logs.Error("Error creating StderrPipe for Cmd")
		logs.Error(cerr)
	}
	errscanner := bufio.NewScanner(cmderrReader)
	go func() {
		for errscanner.Scan() {
			logs.Info(errscanner.Text())
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
func Md5V2(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
//handle error
func Handleerror(errors error)(error){
	logs.Error(errors)
	return errors
}