//==================================
//  * Name：DataSync
//  * DateTime：2019/08/15
//  * File: db.go
//  * Note: db common config .
//==================================

package db

import (
	"log"

	"database/sql"
)

const (
	MYSQL       = "mysql"
	POSTGRESQL  = "postgres"
	ORACLE      = "oracle"
	SQLITE3     = "sqlite3"
	MSSQLSERVER = "sqlserver"
)

const (
	INSERT = "insert"
	UPDATE = "update"
)

type DBServer struct {
	DBtype      string
	DBDriveName string
	User        string
	Passwd      string
	Host        string
	DBName      string
	DB          *sql.DB
	TablesName  []string
	TableKeys   map[string][]string
}

type QueryInfo struct {
	DBName       string
	TableName    string
	PK           string
	KeyArr       []string
	ConditionArr []string
	Content      map[string]string
}

type ExecInfo struct {
	DBName    string
	TableName string
	Handle    string
	PK        string
	Content   []map[string]string
}

type ReNameMapInfo struct {
	SrcTableName string
	DesTableName string
	NameMap      map[string]string
}

const DBConnMax = 8
const MAXSQLCOUNT = 1000
const MAXPAGECOUNT = 20000

func init() {

}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkErrPanic(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
