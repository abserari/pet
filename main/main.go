package main

import (
	"database/sql"
	admin "github.com/abserari/pet/admin/controller"
	permission "github.com/abserari/pet/permission/controller/gin"
	smservice "github.com/abserari/pet/smservice/controller/gin"
	service "github.com/abserari/pet/smservice/service"
	upload "github.com/abserari/pet/upload/controller/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type funcv struct{}

func (v funcv) OnVerifySucceed(targetID, mobile string) {}
func (v funcv) OnVerifyFailed(targetID, mobile string)  {}

func main() {
	var v funcv

	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(192.168.0.253:3307)/test")
	if err != nil {
		panic(err)
	}

	con := &service.Config{
		Host:           "https://fesms.market.alicloudapi.com/sms/",
		Appcode:        "6f37345cad574f408bff3ede627f7014",
		Digits:         6,
		ResendInterval: 60,
		OnCheck:        v,
	}
	smserviceCon := smservice.New(dbConn, con)
	smserviceCon.RegisterRouter(router.Group("/api/v1/message"))

	adminCon := admin.New(dbConn)
	adminCon.RegisterRouter(router.Group("/api/v1/admin"))
	// start to add token on every API after admin.RegisterRouter
	router.Use(adminCon.JWT.MiddlewareFunc())
	// start to check the user active every time.
	router.Use(adminCon.CheckActive())

	permissionCon := permission.New(dbConn, adminCon.GetID)
	permissionCon.RegisterRouter(router.Group("/api/v1/permission"))

	uploadCon := upload.New(dbConn, "http://0.0.0.1:9573", adminCon.GetID)
	uploadCon.RegisterRouter(router.Group("/api/v1/user"))

	log.Fatal(router.Run(":8000"))
}
