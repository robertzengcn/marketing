package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
//	"net/mail"
)

func Init() {
	beego.AddFuncMap("ValidEmail", ValidEmail)
	beego.AddFuncMap("Contains", Contains)
	//add new function in here
}
type NumberStr interface {
    int64 | float64| string
}

//valid email valid
func ValidEmail(email string) bool {
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
func ContainsType[V NumberStr](s []V, str V)bool{
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

///run command
func Runcommand(cmdName string, cmdArgs ...string) error {
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
	cmderrReader, cerr := cmd.StderrPipe()
	if cerr != nil {
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
func Handleerror(errors error) error {
	logs.Error(errors)
	return errors
}

///io read file return file content
func ReadFile(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	return ioutil.ReadAll(jsonFile)

}

///get top domain from url
func Gettopdomain(urls string) (string, error) {
	url, err := url.Parse(urls)
	if err != nil {
		return "", err
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname, nil
}
func PanicFunc(errorObj error) {
	panic(errorObj.Error())
}
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

///read data from csv file
func Csvfilehandle(filepath string) ([][]string, error) {
	// open file
	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, err
}

//check string is url
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//write data to file, create file if not exist
func Writetofile(filename string, data string) error {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}

