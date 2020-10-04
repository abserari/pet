/*
 * Revision History:
 *     Initial: 2020/10/03 Abserari
 */

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Pet -
type Pet struct {
	PetID     int
	Name      string
	ImagePath string
	EventPath string
	StartDate string
	EndDate   string
}

const (
	mysqlPetCreateTable = iota
	mysqlPetInsert
	mysqlPetLisitValidPet
	mysqlPetInfoByID
	mysqlPetDeleteByID
)

var (
	errInvalidInsert = errors.New("insert pet:insert affected 0 rows")

	petSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			 petId    INT NOT NULL AUTO_INCREMENT,
			 name        VARCHAR(512) UNIQUE DEFAULT NULL,
			 imagePath   VARCHAR(512) DEFAULT NULL,
			 eventPath   VARCHAR(512) DEFAULT NULL,
			 startDate   DATETIME NOT NULL,
			 endDate     DATETIME NOT NULL,
			 PRIMARY KEY (petId)
		 )ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO  %s (name,imagePath,eventPath,startDate,endDate) VALUES (?,?,?,?,?)`,
		`SELECT * FROM %s WHERE unix_timestamp(startDate) <= ? AND unix_timestamp(endDate) >= ? LOCK IN SHARE MODE`,
		`SELECT * FROM %s WHERE petid = ? LIMIT 1 LOCK IN SHARE MODE`,
		`DELETE FROM %s WHERE petid = ? LIMIT 1`,
	}
)

// CreateTable -
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(petSQLString[mysqlPetCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

// InsertPet return  id
func InsertPet(db *sql.DB, tableName string, name string, imagePath string, eventPath string, startDate time.Time, endDate time.Time) (int, error) {
	sql := fmt.Sprintf(petSQLString[mysqlPetInsert], tableName)
	result, err := db.Exec(sql, name, imagePath, eventPath, startDate, endDate)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	petID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(petID), nil
}

// LisitValidPetByUnixDate return pet list which have valid date
func LisitValidPetByUnixDate(db *sql.DB, tableName string, unixtime int64) ([]*Pet, error) {
	var (
		bans []*Pet

		petID     int
		name      string
		imagePath string
		eventPath string
		startDate string
		endDate   string
	)

	sql := fmt.Sprintf(petSQLString[mysqlPetLisitValidPet], tableName)
	rows, err := db.Query(sql, unixtime, unixtime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&petID, &name, &imagePath, &eventPath, &startDate, &endDate); err != nil {
			return nil, err
		}

		ban := &Pet{
			PetID:     petID,
			Name:      name,
			ImagePath: imagePath,
			EventPath: eventPath,
			StartDate: startDate,
			EndDate:   endDate,
		}

		bans = append(bans, ban)
	}

	return bans, nil
}

// InfoByID squery by id
func InfoByID(db *sql.DB, tableName string, id int) (*Pet, error) {
	var ban Pet

	sql := fmt.Sprintf(petSQLString[mysqlPetInfoByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ban.PetID, &ban.Name, &ban.ImagePath, &ban.EventPath, &ban.StartDate, &ban.EndDate); err != nil {
			return nil, err
		}
	}

	return &ban, nil
}

// DeleteByID delete by id
func DeleteByID(db *sql.DB, tableName string, id int) error {
	sql := fmt.Sprintf(petSQLString[mysqlPetDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	return err
}
