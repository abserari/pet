/*
 * Revision History:
 *     Initial: 2020/10/15      oiar
 */

package controller

import (
	"database/sql"
	"github.com/abserari/pet/shop/model/mysql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// ShopController -
type ShopController struct {
	db        *sql.DB
	tableName string
}

// New -
func New(db *sql.DB, tableName string) *ShopController {
	return &ShopController{
		db:        db,
		tableName: tableName,
	}
}

// RegisterRouter -
func (b *ShopController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(b.db, b.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", b.create)
	r.POST("/delete", b.deleteByID)
	r.POST("/info/id", b.infoByID)
	r.POST("/list/adminId", b.listShopByAdminID)
}

func (b *ShopController) create(c *gin.Context) {
	var (
		req struct {
			ShopName  string    `json:"shopName"      binding:"required"`
			Address   string    `json:"address"       binding:"required"`
			Cover 	  string    `json:"cover"         binding:"required"`
			Article   string    `json:"article"       binding:"required"`
			Like      bool      `json:"like"          binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertShop(b.db, b.tableName, req.ShopName, req.Address, req.Cover, req.Article, req.Like)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

func (b *ShopController) listShopByAdminID(c *gin.Context) {
	var (
		req struct {
			ShopID uint64 `json:"shopId"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	Shops, err := mysql.ListShop(b.db, b.tableName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Shops": Shops})
}

func (b *ShopController) infoByID(c *gin.Context) {
	var (
		req struct {
			ShopID uint64 `json:"shopId"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	shop, err := mysql.InfoByID(b.db, b.tableName, req.ShopID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": shop})
}

func (b *ShopController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ShopID int `json:"shopId"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteByID(b.db, b.tableName, req.ShopID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
