package main

import (
	"database/sql"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	admin "github.com/abserari/pet/admin/controller/gin"
	banner "github.com/abserari/pet/banner/controller/gin"
	category "github.com/abserari/pet/category/controller/gin"
	order "github.com/abserari/pet/order/controller/gin"
	permission "github.com/abserari/pet/permission/controller/gin"
	smservice "github.com/abserari/pet/smservice/controller/gin"
	service "github.com/abserari/pet/smservice/service"
	upload "github.com/abserari/pet/upload/controller/gin"
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

	adminCon := admin.New(dbConn)
	bannerCon := banner.New(dbConn, "schedule")
	smserviceCon := smservice.New(dbConn, con)

	bannerCon.RegisterRouter(router.Group("/api/v1/schedule"))
	smserviceCon.RegisterRouter(router.Group("/api/v1/message"))

	order.Register(router, dbConn)
	category.Register(dbConn, "students", "test", router)

	JWTMiddleware = &jwt.GinJWTMiddleware{
		Realm:   "Template",
		Key:     []byte("hydra"),
		Timeout: 24 * time.Hour,
	}

	getUID := adminCon.ExtendJWTMiddleWare(JWTMiddleware)

	router.POST("/api/v1/admin/login", JWTMiddleware.LoginHandler)

	router.Use(func(c *gin.Context) {
		JWTMiddleware.MiddlewareFunc()(c)
	})

	router.Use(admin.CheckActive(adminCon, getUID))

	permissionCon := permission.New(dbConn)
	router.Use(permission.CheckPermission(permissionCon, getUID))
	permissionCon.RegisterRouter(router.Group("/api/v1/permission"))

	adminCon.RegisterRouter(router.Group("/api/v1/admin"))

	uploadCon := upload.New(dbConn, "http://0.0.0.1:9573", getUID)
	uploadCon.RegisterRouter(router.Group("/api/v1/user"))

	router.Run(":8000")
}
