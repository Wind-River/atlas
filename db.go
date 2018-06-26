package main

// Licensing: Apache-2.0
/*
 *  Copyright (c) 2017 Wind River Systems, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
 * OR CONDITIONS OF ANY KIND, either express or implied.
 */

/**** https://github.com/mattn/go-sqlite3
 Copyright (c) 2014 Yasuhiro Matsumoto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
****/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// File state globals
var theDB *sql.DB

const CONFIG_RECORD = "Config_Info"

type Supplier_struct struct {
	Id         string `json:"id",omitempty`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	SKU_Symbol string `json:"sku_symbol"` // WR, INTEL, SAMSG
	Short_Id   string `json:"short_id"`   // WR-2e72b7
	Passwd     string `json:"passwd"`
	Type       string `json:"type"` // commercial, non-profit, individual
	Url        string `json:"url"`
}

type ledgerNode struct {
	Id          string `json:"id"`          // Primary key Id
	Name        string `json:"name"`        // Fullname
	ShortId     string `json:"short_id"`    //	1-5 alphanumeric characters (unique)
	IPAddress   string `json:"ip_address"`  // IP address - e.g., 147.11.153.122
	Port        int    `json:"port"`        // Port e.g., 5000
	UUID        string `json:"uuid"`        // 	UUID provide w/previous registration
	Label       string `json:"label"`       // 1-5 words display description
	Description string `json:"description"` // 2-3 sentence description
	Available   int    `json:"available"`   // 0 or 1 int (boolean)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Open the database.
func openDB() {

	var err error
	// Using SQLite
	theDB, err = sql.Open("sqlite3", MAIN_config.DatabaseFile)
	if err != nil {
		panic(err)
	}
	if theDB == nil {
		panic("DB nil")
	}
}

// Initialize the database
func InitializeDB() {

	fmt.Println("Initializing DB ...")
	createDBTables()

	fmt.Println()
	jsonData, _ := dumpDBTable("Applications")
	fmt.Println(jsonData)
	fmt.Println()
}

// Create database tables
func createDBTables() {

	// TODO: create and store database model version

	openDB()
	defer theDB.Close()

	// Create the Ledger Node Table - a list of Ledger node records
	sql_cmd := `
	CREATE TABLE IF NOT EXISTS NetworkSpaceList (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT,
		Description TEXT,
		Status TEXT,
		PublicKey TEXT,
		PrivateKey TEXT,
		InsertedDatetime DATETIME
	);
	`
	_, err := theDB.Exec(sql_cmd)
	if err != nil {
		panic(err)
	}

	// Set Name field to be unique in the Network table. If Insert detects conflict
	// on Name it will replace existing with new record.
	sql_cmd = `CREATE UNIQUE INDEX idx_Name 
				ON NetworkSpaceList (Name);`

	_, err = theDB.Exec(sql_cmd)
	if err != nil {
		fmt.Println(err)
	}

	// Create the Ledger Node Table - a list of Ledger node records
	sql_cmd = `
	CREATE TABLE IF NOT EXISTS LedgerNodes (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		NetworkSpaceID	INTEGER,
        UUID TEXT,
		Name TEXT,
		Alias TEXT,
		API_URL TEXT,
		Description TEXT,
		Status TEXT,
		PublicKey TEXT,
		InsertedDatetime DATETIME
	);
	`
	_, err = theDB.Exec(sql_cmd)
	if err != nil {
		panic(err)
	}

	// Set UUID field to be unique in the Ledger Node table. If Insert detects conflict
	// on UUID it will replace existing with new record.
	sql_cmd = `CREATE UNIQUE INDEX idx_Ledgers_UUID 
				ON LedgerNodes (UUID);`
	_, err = theDB.Exec(sql_cmd)
	if err != nil {
		fmt.Println(err)
	}
}

func dumpDBTable(table_name string) (string, error) {

	defer theDB.Close()
	// Prepare statement to get the native types.
	stmt, err := theDB.Prepare(fmt.Sprintf("SELECT * FROM %s", table_name))

	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}
		tableData = append(tableData, entry)
	}

	jsonData99, err := json.Marshal(tableData)

	return string(jsonData99), nil
}

// A boolean function that determines if a record with uuid exits in the db/
func ApplicationExists(uuid string) bool {

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT UUID FROM Applications WHERE UUID=?", uuid)

	checkErr(err)

	// rows.Next () is a boolean that says whether a next record exists
	// If there is a next record (i.e., at least one) then true else false
	if rows.Next() {
		rows.Close()
		return true
	} else {
		rows.Close()
		return false
	}
}

// A boolean function that determines if a ledger node record with uuid exits in the db
func LedgerNodeExists(uuid string) bool {

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT UUID FROM LedgerNodes WHERE UUID=?", uuid)

	checkErr(err)

	// rows.Next () is a boolean that says whether a next record exists
	// If there is a next record (i.e., at least one) then true else false
	if rows.Next() {
		rows.Close()
		return true
	} else {
		rows.Close()
		return false
	}
}

// Get the Ledger node info
func GetLedgerNodeInfo(uuid string, lnode *ledgerNode) {

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT UUID, Name, Short_Id, IP_Address, Port, Label, Description FROM LedgerNodes WHERE UUID=?", uuid)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&lnode.UUID, &lnode.Name, &lnode.ShortId, &lnode.IPAddress, &lnode.Port,
			&lnode.Label, &lnode.Description)
		checkErr(err)
	}
	rows.Close()
}

// Returns the db ID for Network record. Returns 0 if does not exist. Returns -1 if error encoutered
func GetNetworkSpaceIDFromDB(networkName string) (Id int, err error) {
	var networkSpaceID int

	openDB()
	defer theDB.Close() // Query for  network space id
	rows, err := theDB.Query("SELECT Id FROM NetworkSpaceList WHERE Name=?", strings.ToLower(networkName))
	if err != nil {
		return -1, err
	}

	// rows.Next () is a boolean that says whether a next record exists
	if !rows.Next() {
		// No record exists. Return 0 for none.
		return 0, nil
	}

	// ID exists return ID.
	err = rows.Scan(&networkSpaceID)
	rows.Close()
	if err != nil {
		return -1, err
	}
	return networkSpaceID, nil
}

// Insert network space record into the DB
func AddNetworkSpaceToDB(record NetworkSpaceRecord) error {

	// Check if network already exists
	networkSpaceID, err := GetNetworkSpaceIDFromDB(record.Name)
	if networkSpaceID > 0 {
		return fmt.Errorf("Network '%s' already exists", record.Name)
	}
	if err != nil {
		return err
	}

	// Network does not already exist. Procceed.
	openDB()
	defer theDB.Close()

	sql_additem := `
	INSERT OR REPLACE INTO NetworkSpaceList (
		Name,
		Description,
		Status,
		PublicKey,
		InsertedDatetime
		) values(?, ?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(strings.ToLower(record.Name), record.Description, record.Status, record.PublicKey)
	if err != nil {
		return fmt.Errorf("error - inserting record for network: %s into db", record.Name)
	}

	// successfully inserted. No errors.
	return nil
}

func GetLedgerNodeListDB(networkName string) ([]LedgerNodeRecord, error) {
	var list []LedgerNodeRecord
	var record LedgerNodeRecord

	networkSpaceID, err := GetNetworkSpaceIDFromDB(networkName)
	if networkSpaceID == 0 {
		return nil, fmt.Errorf("Network name space '%s' does not exist", networkName)
	}
	if err != nil {
		return nil, err
	}

	openDB()
	defer theDB.Close()

	/****
	// Query for  network space id
	rows, err := theDB.Query("SELECT Id FROM NetworkSpaceList WHERE Name=?", strings.ToLower(networkName))
	checkErr(err)

	// rows.Next () is a boolean that says whether a next record exists
	// If there is a next record (i.e., at least one) then true else false
	var networkSpaceID int
	if rows.Next() {
		err = rows.Scan(&networkSpaceID)
		checkErr(err)
		rows.Close()
	} else {
		// Network not found - error
		// TODO: handle error better - tell client no supplier.
		rows.Close()
		return nil, fmt.Errorf("Network Name Space not found: %s", networkName)
	}
	******/

	// The network exists, proceed.
	fmt.Println("network id:", networkSpaceID)
	rows, err := theDB.Query("SELECT Name, Alias, UUID, Description, API_URL, Status, PublicKey, InsertedDatetime FROM LedgerNodes WHERE NetworkSpaceID=?", networkSpaceID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error accessing record for network: %s", networkName)
	}

	// initialize the list to the empty list.
	list = make([]LedgerNodeRecord, 0)
	for rows.Next() {
		err = rows.Scan(&record.Name, &record.Alias, &record.UUID, &record.Description, &record.APIURL, &record.Status, &record.PublicKey, &record.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("error accessing ledger node records")
		}
		list = append(list, record)
	}
	rows.Close()
	return list, nil
}

func GetNetworkSpaceListDB() []NetworkSpaceRecord {

	var list []NetworkSpaceRecord
	var record NetworkSpaceRecord

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT Name, Description, Status, InsertedDatetime FROM NetworkSpaceList")
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&record.Name, &record.Description, &record.Status, &record.Timestamp)
		checkErr(err)
		list = append(list, record)
	}
	rows.Close()
	return list
}

// Insert Ledger node record into the DB
func AddLedgerNodeToDB(record LedgerNodeRecord) error {

	networkSpaceID, err := GetNetworkSpaceIDFromDB(record.NetworkName)
	if networkSpaceID == 0 {
		// successful accessed db but network not found.
		return fmt.Errorf("Network name space '%s' does not exist", record.NetworkName)
	}
	if err != nil {
		// error occurred accessing db
		fmt.Println(err)
		return fmt.Errorf("system error - could not access data store")
	}

	openDB()
	defer theDB.Close()

	// TODO: Check Node UUID is properly formatted.

	/****************************
	// Query for  network space id
	rows, err := theDB.Query("SELECT Id FROM NetworkSpaceList WHERE Name=?", strings.ToLower(node_record.NetworkName))
	checkErr(err)

	// rows.Next () is a boolean that says whether a next record exists
	// If there is a next record (i.e., at least one) then true else false
	var networkSpaceID int
	if rows.Next() {
		err = rows.Scan(&networkSpaceID)
		checkErr(err)
		rows.Close()
	} else {
		// Network not found - error
		// TODO: handle error better - tell client no supplier.
		fmt.Println("Network Space not found:", node_record.NetworkName)
		rows.Close()
		return err
	}
	*****************/

	// TODO: ping node to obtain status

	// TODO:
	// For an upadate UUID, Network Name and Public key decrypt need to match in order to update.
	// The network name would be decrypted using the public key and hence match.

	sql_additem := `
	INSERT OR REPLACE INTO LedgerNodes (
		NetworkSpaceID,
		UUID, 
		Name, 
		Alias,
		API_URL,
		Description,
		Status,
		PublicKey,
		InsertedDatetime
		) values(?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(networkSpaceID, record.UUID, record.Name, record.Alias, record.APIURL,
		record.Description, record.Status, record.PublicKey)
	if err != nil {
		panic(err)
	}

	return nil
}

// Get ledger node record
func GetLedgerNodesFromDB() []LedgerNode {

	var list []LedgerNode
	var node LedgerNode

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT UUID, Name, Short_Id, API_URL, Label, Description FROM LedgerNodes")
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&node.UUID, &node.Name, &node.ShortId, &node.API_Address,
			&node.Label, &node.Description)
		checkErr(err)
		list = append(list, node)
	}
	rows.Close() //good habit to close
	return list
}

// Get most recently reported Ledger API network address
func GetLedgerAPIAddress(ip_address *string, port *int) {

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT Ledger_IP_Addr, Ledger_Port FROM SystemConfig WHERE Config_Name=?", CONFIG_RECORD)
	checkErr(err)
	rows.Next()

	err = rows.Scan(ip_address, port)
	checkErr(err)
	fmt.Println("Ledger IP       = ", *ip_address)
	fmt.Println("Ledger Port     =", *port)
	rows.Close()
}

// Update the Ledger API address
func UpdateLedgerAPIAddress(ip_address string, port int) {
	openDB()
	defer theDB.Close()

	stmt, err := theDB.Prepare("UPDATE SystemConfig SET  Ledger_IP_Addr=?, Ledger_Port=?  WHERE Config_Name=?")
	checkErr(err)
	_, err = stmt.Exec(ip_address, port, CONFIG_RECORD)
}

// Add supplier record to the database
func AddSupplierToDB(s Supplier_struct) {

	openDB()
	defer theDB.Close()
	sql_additem := `
	INSERT OR REPLACE INTO Suppliers (
		UUID, 
		Name, 
		SKU_Symbol, 
		Short_Id, 
		Passwd, 
		Type, 
		Url, 
		InsertedDatetime
		) values(?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	res, err2 := stmt.Exec(s.UUID, s.Name, s.SKU_Symbol, s.Short_Id, s.Passwd, s.Type, s.Url)
	if err2 != nil {
		panic(err2)
	}

	// for debugging - could probably be removed.
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error:")
		fmt.Println("last Id inserted", id)
		fmt.Println("Inserted")
		fmt.Println(s)
	}
}

// Get Supplier record from database.
func GetSupplier(db *sql.DB) []Supplier_struct {
	sql_readall := `
	SELECT Id, UUID, Name, SKU_Symbol, Short_Id, Passwd, Type, Url FROM Suppliers
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Supplier_struct
	for rows.Next() {
		item := Supplier_struct{}
		err2 := rows.Scan(&item.Id, &item.UUID, &item.Name, &item.SKU_Symbol, &item.Short_Id, &item.Passwd, &item.Type, &item.Url)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}
