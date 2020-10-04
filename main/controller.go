/*
 * Revision History:
 *     Initial: 2020/10/03 Abserari
 */

package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PetController -
type PetController struct {
	db        *sql.DB
	tableName string
}

// New -
func New(db *sql.DB, tableName string) *PetController {
	return &PetController{
		db:        db,
		tableName: tableName,
	}
}

// RegisterRouter -
func (b *PetController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := CreateTable(b.db, b.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", b.create)
	r.POST("/delete", b.deleteByID)
	r.POST("/info/id", b.infoByID)
	r.POST("/list/date", b.lisitValidPetByUnixDate)
}

func (b *PetController) create(c *gin.Context) {
	var (
		req struct {
			Name      string    `json:"name"      binding:"required"`
			ImagePath string    `json:"imageurl"  binding:"required"`
			EventPath string    `json:"eventurl"  binding:"required"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := InsertPet(b.db, b.tableName, req.Name, req.ImagePath, req.EventPath, req.StartDate, req.EndDate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

func (b *PetController) lisitValidPetByUnixDate(c *gin.Context) {
	var (
		req struct {
			Unixtime int64 `json:"unixtime"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	pets, err := LisitValidPetByUnixDate(b.db, b.tableName, req.Unixtime)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "pets": pets})
}

func (b *PetController) infoByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := InfoByID(b.db, b.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}

func (b *PetController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = DeleteByID(b.db, b.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
