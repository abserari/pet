/*
 * Revision History:
 *     Initial: 2020/10/15      oiar
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

type Shop struct {
	ShopID     uint64
	ShopName   string
	Address    string
	Cover      string
	Article    string
	Like       bool
}

const (
	mysqlShopCreateTable = iota
	mysqlShopInsert
	mysqlShopList
	mysqlShopInfoByID
	mysqlShopDeleteByID
)

var (
	errInvalidInsert = errors.New("insert shop:insert affected 0 rows")

	shopSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			shopId      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			shopName    VARCHAR(512) UNIQUE DEFAULT NULL DEFAULT ' ',
			address		VARCHAR(512) NOT NULL DEFAULT ' ',
			cover       VARCHAR(512) NOT NULL DEFAULT ' ',
			article     VARCHAR(512) NOT NULL DEFAULT ' ',
			like		BOOLEAN NOT NULL DEFAULT FALSE,
			PRIMARY KEY (shopId)
		)ENGINE=InnoDB AUTO_INCREMENT=1000000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s (shopName,address,cover,article,like) VALUES (?,?,?,?,?)`,
		`SELECT * FROM %s`,
		`SELECT * FROM %s WHERE shopId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`DELETE FROM %s WHERE shopId = ? LIMIT 1`,
	}
)

// CreateTable -
func CreateTable(db *sql.DB, tableName string) error {
	s := fmt.Sprintf(shopSQLString[mysqlShopCreateTable], tableName)
	_, err := db.Exec(s)
	return err
}

// InsertShop return id
func InsertShop(db *sql.DB, tableName, shopName, address, cover, article string, like bool) (int, error) {
	s := fmt.Sprintf(shopSQLString[mysqlShopInsert], tableName)
	result, err := db.Exec(s, shopName, address, cover, article, like)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	shopID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(shopID), nil
}

// listShop return shop list
func ListShop(db *sql.DB, tableName string) ([]*Shop, error) {
	var (
		shops []*Shop

		shopID   uint64
		shopName string
		address  string
		cover    string
		article  string
		like     bool
	)

	s := fmt.Sprintf(shopSQLString[mysqlShopList], tableName)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&shopID, &shopName, &address, &cover, &article, &like); err != nil {
			return nil, err
		}

		shop := &Shop{
			ShopID:   shopID,
			ShopName: shopName,
			Address:  address,
			Cover:    cover,
			Article:  article,
			Like:     like,
		}

		shops = append(shops, shop)
	}

	return shops, nil
}

// InfoByID query by id
func InfoByID(db *sql.DB, tableName string, shopId uint64) (*Shop, error) {
	var shop Shop

	s := fmt.Sprintf(shopSQLString[mysqlShopInfoByID], tableName)
	rows, err := db.Query(s, shopId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&shop.ShopID, &shop.ShopName, &shop.Address, &shop.Cover, &shop.Article, &shop.Like); err != nil {
			return nil, err
		}
	}

	return &shop, nil
}

// DeleteByID delete by id
func DeleteByID(db *sql.DB, tableName string, shopId int) error {
	s := fmt.Sprintf(shopSQLString[mysqlShopDeleteByID], tableName)
	_, err := db.Exec(s, shopId)
	return err
}
