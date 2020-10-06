/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	mysql "github.com/abserari/pet/admin/model/mysql"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	errActive          = errors.New("Admin active is wrong")
	errUserIDNotExists = errors.New("Get Admin ID is not exists")
	errUserIDNotValid  = func(value interface{}) error {
		return errors.New(fmt.Sprintf("Get Admin ID is not valid. Is %s", value))
	}
)

// Controller external service interface
type Controller struct {
	db *sql.DB
}

// New create an external service interface
func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

// RegisterRouter register router. It fatal because there is no service if register failed.
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}
	// the jwt middleware

	name := "Admin"
	password := "111111"
	err := mysql.CreateTable(c.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}
	jwtMiddleware, err := c.newJWTMiddleware()
	if err != nil {
		log.Fatal(err)
	}

	// login and refresh token.
	r.POST("/login")
	r.GET("/refresh_token", jwtMiddleware.RefreshHandler)

	// start to add token on every API after admin.RegisterRouter
	r.Use(jwtMiddleware.MiddlewareFunc())
	// start to check the user active every time.
	r.Use(c.checkActive())

	// admin crud API
	r.POST("/create", c.create)
	r.POST("/modify/email", c.modifyEmail)
	r.POST("/modify/mobile", c.modifyMobile)
	r.POST("/modify/password", c.modifyPassword)
	r.POST("/modify/active", c.modifyAdminActive)
}

func (c *Controller) GetID(ctx *gin.Context) (uint32, error) {
	id, ok := ctx.Get("userID")
	if !ok {
		return 0, errUserIDNotExists
	}

	v, ok := id.(uint32)
	if !ok {
		return 0, errUserIDNotValid(id)
	}
	return v, nil
}

//CheckActive middleware that checks the active
func (c *Controller) checkActive() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		a, err := c.GetID(ctx)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		active, err := mysql.IsActive(c.db, a)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		if !active {
			_ = ctx.AbortWithError(http.StatusLocked, errActive)
			ctx.JSON(http.StatusLocked, gin.H{"status": http.StatusLocked})
			return
		}
	}
}

func (con *Controller) newJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test-pet",
		Key:         []byte("moli-tech-cats-member"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		// use data as userID here.
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{
				"userID": data,
			}
		},
		// just get the ID
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["userID"]
		},
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			return con.Login(ctx)
		},
		// no need to check user valid every time.
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
func (c *Controller) create(ctx *gin.Context) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	//Default password
	if admin.Password == "" {
		admin.Password = "111111"
	}

	err = mysql.Create(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"    binding:"required"`
			Email   string `json:"email"       binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyEmail(c.db, admin.AdminID, &admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyMobile(ctx *gin.Context) {
	var (
		admin struct {
			AdminID uint32 `json:"admin_id"     binding:"required"`
			Mobile  string `json:"mobile"       binding:"required,numeric,len=11"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMobile(c.db, admin.AdminID, &admin.Mobile)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyPassword(ctx *gin.Context) {
	var (
		admin struct {
			AdminID     uint32 `json:"admin_id"       binding:"required"`
			Password    string `json:"password"       binding:"printascii,min=6,max=30"`
			NewPassword string `json:"new_password"   binding:"printascii,min=6,max=30"`
			Confirm     string `json:"confirm"        binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.NewPassword == admin.Password {
		ctx.Error(err)
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	if admin.NewPassword != admin.Confirm {
		ctx.Error(err)
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.ModifyPassword(c.db, admin.AdminID, &admin.Password, &admin.NewPassword)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyAdminActive(ctx *gin.Context) {
	var (
		admin struct {
			CheckID     uint32 `json:"check_id"    binding:"required"`
			CheckActive bool   `json:"check_active"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyAdminActive(c.db, admin.CheckID, admin.CheckActive)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

//Login JWT validation
func (c *Controller) Login(ctx *gin.Context) (uint32, error) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		return 0, err
	}

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Password)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
