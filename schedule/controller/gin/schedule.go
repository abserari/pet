/*
 * Revision History:
 *     Initial: 2020/10/17       oiar
 */

package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/abserari/pet/schedule/model/mysql"
	"github.com/gin-gonic/gin"
)

// ScheduleController -
type ScheduleController struct {
	db        *sql.DB
	tableName string
	getUID  func(c *gin.Context) (uint64, error)
}

// New -
func New(db *sql.DB, tableName string, getUID func(c *gin.Context) (uint64, error)) *ScheduleController {
	return &ScheduleController{
		db:        db,
		tableName: tableName,
		getUID: getUID,
	}
}

// RegisterRouter -
func (b *ScheduleController) RegisterRouter(r gin.IRouter) {
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
	r.POST("/list/adminId", b.listScheduleByAdminID)
}

func (b *ScheduleController) create(c *gin.Context) {
	var (
		req struct {
			Date      string    `json:"date"     binding:"required"`
			Time      string    `json:"time"     binding:"required"`
			Note      string    `json:"note"     binding:"required"`
		}
	)

	adminId, err := b.getUID(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertSchedule(b.db, b.tableName, adminId, req.Date, req.Time, req.Note)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

func (b *ScheduleController) listScheduleByAdminID(c *gin.Context) {
	var (
		req struct {
			AdminID uint64 `json:"adminId"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	Schedules, err := mysql.ListValidScheduleByAdminID(b.db, b.tableName, req.AdminID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Schedules": Schedules})
}

func (b *ScheduleController) infoByID(c *gin.Context) {
	var (
		req struct {
			ScheduleID uint64 `json:"scheduleId"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := mysql.InfoByID(b.db, b.tableName, req.ScheduleID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}

func (b *ScheduleController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ScheduleID uint64 `json:"scheduleId"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteByID(b.db, b.tableName, req.ScheduleID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
