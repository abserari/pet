/*
 * Revision History:
 *     Initial: 2020/10/05        Abserari
 */

package controller

import (
	"errors"
	"net/http"

	permission "github.com/abserari/pet/permission/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	//FirstURL for check whether they meet the requirements to execution middleware
	FirstURL      = "/api/v1/permission/addurl"
	errPermission = errors.New("Admin permission is wrong")
)

//CheckPermission middleware that checks the permission
func (c *Controller)CheckPermission() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		var check = true
		URLL := ctx.Request.URL.Path

		adminID, err := c.getIDFunc(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		adRole, err := permission.AdminGetRoleMap(c.db, adminID)
		if err != nil {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		urlRole, err := permission.URLPermissions(c.db, &URLL)
		if err != nil {
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		urlroleid, err := permission.URLPermissions(c.db, &FirstURL)
		if err != nil {
			ctx.AbortWithError(http.StatusVariantAlsoNegotiates, err)
			return
		}

		for roleid := range urlroleid {
			for adminid := range adRole {
				if roleid == adminid {
					check = false
				}
			}
		}

		for urlkey := range urlRole {
			for adkey := range adRole {
				if urlkey == adkey {
					check = true
				}
			}
		}

		if !check {
			ctx.AbortWithError(http.StatusForbidden, errPermission)
			ctx.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		}

	}
}
