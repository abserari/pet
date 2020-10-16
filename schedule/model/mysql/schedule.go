/*
 * Revision History:
 *     Initial: 2020/10/16       oiar
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

// Schedule -
type Schedule struct {
	ScheduleID uint64
	AdminID    uint64
	Date       string
	Time       string
	Note       string
}

const (
	mysqlScheduleCreateTable = iota
	mysqlScheduleInsert
	mysqlScheduleListValidSchedule
	mysqlScheduleInfoByID
	mysqlScheduleDeleteByID
)

var (
	errInvalidInsert = errors.New("insert schedule:insert affected 0 rows")

	ScheduleSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			scheduleId  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			adminId     BIGINT UNSIGNED NOT NULL,
			date        VARCHAR(512) UNIQUE DEFAULT NULL,
			time 		VARCHAR(512) UNIQUE DEFAULT NULL,
			note        VARCHAR(512) UNIQUE DEFAULT NULL,
			PRIMARY KEY (scheduleId),
			KEY (adminId)
		)ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s (adminId, date,time,note) VALUES (?,?,?,?)`,
		`SELECT * FROM %s WHERE adminId = ? LOCK IN SHARE MODE`,
		`DELETE FROM %s WHERE scheduleId = ? LIMIT 1`,
	}
)

// CreateTable -
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(ScheduleSQLString[mysqlScheduleCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

// InsertSchedule return id
func InsertSchedule(db *sql.DB, tableName string, adminId uint64, date string, time string, note string) (int, error) {
	sql := fmt.Sprintf(ScheduleSQLString[mysqlScheduleInsert], tableName)
	result, err := db.Exec(sql, adminId, date, time, note)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	scheduleID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(scheduleID), nil
}

// ListValidScheduleByAdminID return schedule list which have valid date
func ListValidScheduleByAdminID(db *sql.DB, tableName string, adminId uint64) ([]*Schedule, error) {
	var (
		bans []*Schedule

		ScheduleID uint64
		AdminID    uint64
		Date       string
		Time       string
		Note       string
	)

	s := fmt.Sprintf(ScheduleSQLString[mysqlScheduleListValidSchedule], tableName)
	rows, err := db.Query(s, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ScheduleID, &AdminID, &Date, &Time, &Note); err != nil {
			return nil, err
		}

		ban := &Schedule{
			ScheduleID:  ScheduleID,
			AdminID: AdminID,
			Date: Date,
			Time: Time,
			Note: Note,
		}

		bans = append(bans, ban)
	}

	return bans, nil
}

// InfoByID query by id
func InfoByID(db *sql.DB, tableName string, scheduleID uint64) (*Schedule, error) {
	var ban Schedule

	s := fmt.Sprintf(ScheduleSQLString[mysqlScheduleInfoByID], tableName)
	rows, err := db.Query(s, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ban.ScheduleID, &ban.AdminID, &ban.Date, &ban.Time, &ban.Note); err != nil {
			return nil, err
		}
	}

	return &ban, nil
}

// DeleteByID delete by id
func DeleteByID(db *sql.DB, tableName string, scheduleID uint64) error {
	s := fmt.Sprintf(ScheduleSQLString[mysqlScheduleDeleteByID], tableName)
	_, err := db.Exec(s, scheduleID)
	return err
}
