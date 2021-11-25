package utils

import (
    beego "github.com/beego/beego/v2/server/web"
	"regexp"
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