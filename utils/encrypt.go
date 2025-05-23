package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strings"
	//"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
	//"google.golang.org/protobuf/internal/errors"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}



func GetSecret() (string,error) {
	key,err:=beego.AppConfig.String("dataEncrykey")
	// This should be in an env file in production
	if err != nil { 
		return "", err
	}

	keystring:=strings.TrimSpace(key)

	return keystring,nil
	
}
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) ([]byte,error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// panic(err)
		return nil,err
	}
	return data,nil
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text string) (string, error) {
	MySecret, err := GetSecret()
	if(err != nil){
		return "", err
	}
	if(len(MySecret) < 1){
		return "", errors.New("secret key is empty")
	}
	logs.Info(MySecret)
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text string) (string, error) {
	MySecret, err := GetSecret()
	if(err != nil){
		return "", err
	}
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText,derr := Decode(text)
	if(derr != nil){
		return "",derr
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

