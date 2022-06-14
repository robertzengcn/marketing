package tests

import (
	"net/http"
	"net/http/httptest"
	// "marketing/models"
	"testing"
	"fmt"
	"runtime"
	"path/filepath"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
	log "github.com/beego/beego/v2/core/logs"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	"github.com/beego/i18n"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	_ "marketing/routers"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	fmt.Println(apppath)
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.TestBeegoInit(apppath)
	mysql_user, uerr := config.String("mysql_user")
	if(uerr!=nil){
		panic(uerr)
	}
	mysql_pass, perr := config.String("mysql_pass")
	if(perr!=nil){
		panic(perr)	
	}
	mysql_host, hoerr := config.String("mysql_host")
	if(hoerr!=nil){
		panic(hoerr)	
	}
	mysql_port, poerr := config.String("mysql_port")
	if(poerr!=nil){
		panic(poerr)
	}
	mysql_dbname, dbname_err := config.String("mysql_dbname")
	if(dbname_err!=nil){
		panic(dbname_err)
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	 mysqlconurl:=mysql_user+":"+mysql_pass+"@tcp("+mysql_host+":"+mysql_port+")/"+mysql_dbname+"?charset=utf8&parseTime=True&loc=Local"
	
	orm.RegisterDataBase("default", "mysql", mysqlconurl)
	
	// register model
	orm.RunSyncdb("default", false, true)
	orm.Debug,_ = config.Bool("dbdebug")

}

// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost:8080/healthcheck", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	log.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}