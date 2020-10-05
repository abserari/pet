package main

import (
	"database/sql"

	admin "github.com/abserari/pet/admin/controller"
	permission "github.com/abserari/pet/permission/controller/gin"
	smservice "github.com/abserari/pet/smservice/controller/gin"
	service "github.com/abserari/pet/smservice/service"
	upload "github.com/abserari/pet/upload/controller/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// JWTMiddleware should be exported for user authentication.
	JWTMiddleware *jwt.GinJWTMiddleware
)

type funcv struct{}

func (v funcv) OnVerifySucceed(targetID, mobile string) {}
func (v funcv) OnVerifyFailed(targetID, mobile string)  {}

func main() {
	var v funcv

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
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

	permissionCon := permission.New(dbConn)
	router.Use(permission.CheckPermission(permissionCon, adminCon.GetID))
	permissionCon.RegisterRouter(router.Group("/api/v1/permission"))

	uploadCon := upload.New(dbConn, "http://0.0.0.1:9573", adminCon.GetID)
	uploadCon.RegisterRouter(router.Group("/api/v1/user"))

	router.Run(":8000")
}
