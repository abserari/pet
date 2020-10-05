/*
 * Revision History:
 *     Initial: 2019/03/14        Yang ChengKai
 */

package controller

import (
	"errors"
	"net/http"

	ginjwt "github.com/appleboy/gin-jwt"
	mysql "github.com/abserari/pet/admin/model/mysql"
	"github.com/gin-gonic/gin"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

var (
	errActive          = errors.New("Admin active is wrong")
	errUserIDNotExists = errors.New("Get Admin ID is wrong")
)

//ExtendJWTMiddleWare extend JWTMiddleWare
func (c *Controller) ExtendJWTMiddleWare(JWTMiddleware *ginjwt.GinJWTMiddleware) func(ctx *gin.Context) (uint32, error) {
	JWTMiddleware.Authenticator = func(ctx *gin.Context) (interface{}, error) {
		return c.Login(ctx)
	}

	JWTMiddleware.PayloadFunc = func(data interface{}) ginjwt.MapClaims {
		return ginjwt.MapClaims{
			"userID": data,
		}
	}

	JWTMiddleware.IdentityHandler = func(claims jwt.MapClaims) interface{} {
		return claims["userID"]
	}

	return func(ctx *gin.Context) (uint32, error) {
		id, ok := ctx.Get("userID")
		if !ok {
			return 0, errUserIDNotExists
		}

		v := id.(float64)
		return uint32(v), nil
	}
}

//CheckActive middleware that checks the active
func CheckActive(c *Controller, getUID func(ctx *gin.Context) (uint32, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		a, err := getUID(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		active, err := mysql.IsActive(c.db, a)
		if err != nil {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		if !active {
			ctx.AbortWithError(http.StatusLocked, errActive)
			ctx.JSON(http.StatusLocked, gin.H{"status": http.StatusLocked})
			return
		}
	}
}
