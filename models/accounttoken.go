package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"

	// guuid "github.com/google/uuid"
	"github.com/beego/beego/v2/core/config"
	jwt "github.com/golang-jwt/jwt/v5"
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

type UserClaim struct {
	jwt.RegisteredClaims
	Name string
	AccountId    int64 
	Email string
	Roles  []string 
}
///gen account token
func (u *AccountToken) GenAccounttoken(account *Account) (token string,err error){
	//token=guuid.NewString()
	var accountrolearr []string
	// cliams:=jwt.MapClaims{
	// 	"account_id": account.Id,
	// 	"email": account.Email,
	// 	"roles": accountrolearr,
	// 	"nbf": time.Now().Unix(),
	// 	"exp": time.Now().AddDate(0, 0, 2).Unix(),
	// 	"iat": time.Now().Unix(),
	// }
	cliams:=UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		AccountId:account.Id,
		Email: account.Email,
		Name: account.Name,	
		// Roles: accountrolearr,
	}
	if(account.Roles!=nil){
		for _,element:=range account.Roles{
			accountrolearr=append(accountrolearr,element.Name)
		}
		cliams.Roles=accountrolearr
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
	ac.TokenExpired=now.AddDate(0, 1, 0)
	_, err = o.Insert(&ac)
	return token,err
}
//generate token use jwt token
func (u *AccountToken) GenAccounttokenjwt(claims *UserClaim) (string,error){
	token_val, terr := config.String("jwt_token_key")
	if(terr!=nil){
		return "", terr
	}
	token_sec:=[]byte(token_val)
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tokens.SignedString(token_sec)
	if(err!=nil){
		return "", err
	}
	return tokenString,nil
}
///parse token
func (u *AccountToken) ParseAccounttokenjwt(token string,userClaim *UserClaim)(*UserClaim,error){
	token_val, terr := config.String("jwt_token_key")
	if(terr!=nil){
		return nil,terr
	}
	tokenStruct, _ :=jwt.ParseWithClaims(token, userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(token_val), nil
	})
	if !tokenStruct.Valid {
		return nil,jwt.ErrSignatureInvalid
	}
	return userClaim,nil

}

///check account token
func (u *AccountToken) CheckAccounttoken(token string) (*AccountToken,error){
	o := orm.NewOrm()
	accToken := new(AccountToken)
	accounttoken:=AccountToken{}
	now := time.Now()
	nf:=now.Format("2006-01-02 15:04:05")
	err:=o.QueryTable(accToken).Filter("token_val", token).Filter("token_expired__gt", nf).One(&accounttoken)
	if(err!=nil){
		return nil, err
	}else{
		//decode token
		userClaim:=UserClaim{}
		_,perr:=u.ParseAccounttokenjwt(token,&userClaim)
		if(perr!=nil){
			return nil, perr
		}
		if(accounttoken.Account.Id!=userClaim.AccountId){
			return nil, errors.New("token not match account")
		}
		return &accounttoken,nil
	}
}
func (u *AccountToken) DeleteAccountToken(tokenId int64) (int64,error){
	o := orm.NewOrm()
	// accToken := new(AccountToken)
	num, err := o.Delete(&AccountToken{TokenId: tokenId})
	return num,err
}





