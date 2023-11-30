package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	// guuid "github.com/google/uuid"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/beego/beego/v2/core/config"
)
var DefaultAccountToken *AccountToken
type AccountToken struct{
	TokenId int64 `orm:"pk;auto"`
	TokenVal string
	Account   *Account  `orm:"rel(fk);on_delete(do_nothing)"`
	TokenExpired   time.Time `orm:"type(datetime)"`
}



func (u *AccountToken) TableEngine() string {
	return "MYISAM"
} 

func (u *AccountToken) TableName() string {
	return "account_token"
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(AccountToken))
	// create table
	// orm.RunSyncdb("default", false, true)
}
///gen account token
func (u *AccountToken) GenAccounttoken(account *Account) (token string,err error){
	//token=guuid.NewString()
	var accountrolearr []string
	cliams:=jwt.MapClaims{
		"account_id": account.Id,
		"email": account.Email,
		"roles": accountrolearr,
		"nbf": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 0, 2).Unix(),
		"iat": time.Now().Unix(),
	}

	if(account.Roles!=nil){
		for _,element:=range account.Roles{
			accountrolearr=append(accountrolearr,element.Name)
		}
		cliams["roles"]=accountrolearr
	}
	
	token,terr:=u.GenAccounttokenjwt(&cliams)
	if(terr!=nil){
		return "", terr
	}
	o := orm.NewOrm()
	var ac AccountToken
	ac.Account=account
	ac.TokenVal=token
	now := time.Now()
	ac.TokenExpired=now.AddDate(0, 0, 2)
	_, err = o.Insert(&ac)
	return token,err
}
//generate token use jwt token
func (u *AccountToken) GenAccounttokenjwt(claims *jwt.MapClaims) (string,error){
	token_val, terr := config.String("jwt_token_key")
	if(terr!=nil){
		return "", terr
	}
	token_sec:=[]byte(token_val)
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tokens.SignedString(token_sec)
	return tokenString,err
}

///check account token
func (u *AccountToken) CheckAccounttoken(token string) (accounttoken *AccountToken,err error){
	o := orm.NewOrm()
	accToken := new(AccountToken)
	now := time.Now()
	nf:=now.Format("2006-01-02 15:04:05")
	err=o.QueryTable(accToken).Filter("token_val", token).Filter("token_expired__gt", nf).One(&accounttoken)
	if(err!=nil){
		return nil, err
	}else{
		return accounttoken,nil
	}

}





