//==================================
//  * Name：DataSync
//  * DateTime：2019/08/15
//  * File: db.go
//  * Note: db common handle .
//==================================

package db

import (
	"fmt"
	"strconv"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
)

func Setup() {

}

func (db *DBServer) Start() error {
	var err error

	dbtype := db.DBtype
	host := db.Host
	user := db.User
	passwd := db.Passwd
	dbname := db.DBName

	switch dbtype {
	case MYSQL:
		DSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, passwd, host, dbname)
		db.DB, err = sql.Open(MYSQL, DSN)
	case POSTGRESQL:
		DSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, passwd, host, dbname)
		db.DB, err = sql.Open(POSTGRESQL, DSN)
	case ORACLE:
		// [username/[password]@]host[:port][/instance_name]
		DSN := fmt.Sprintf("%s/%s@%s/orcl", user, passwd, host)
		db.DB, err = sql.Open("oci8", DSN)
	case SQLITE3:
		// ./testdb.db
		db.DB, err = sql.Open(SQLITE3, fmt.Sprintf("%s", dbname))
	default:
		panic(fmt.Sprintf("InitDB error . dbtype %s is error .", dbtype))
	}

	if err != nil {
		panic(fmt.Sprintf("InitDB error . error message : %s", err))
	}
	db.DB.SetMaxIdleConns(DBConnMax / 2)
	db.DB.SetMaxOpenConns(DBConnMax)
	return nil
}

func (db *DBServer) Stop() error {
	if db.DB != nil {
		db.DB.Close()
	}
	return nil
}

func (db *DBServer) ReName(resultMap []map[string]string, reNameInfo ReNameMapInfo) []map[string]string {
	res := resultMap
	nameMap := reNameInfo.NameMap
	if nameMap == nil || len(nameMap) == 0 {
		return res
	}
	for i := 0; i < len(res); i++ {
		itemMap := res[i]
		for src, des := range nameMap {
			tmpStr := itemMap[src]
			itemMap[des] = tmpStr
			delete(itemMap, src)
		}
		res[i] = itemMap
	}
	// oracle handle . delete RN .
	if _, ok := res[0]["RN"]; ok {
		for i := 0; i < len(res); i++ {
			itemMap := res[i]
			delete(itemMap, "RN")
		}
	}
	return res
}

func dbQueryString(db *sql.DB, strsql string) ([]map[string]string, error) {
	rows, err := db.Query(strsql)
	fmt.Println("strsql : ", strsql)
	checkErr(err)
	cols, _ := rows.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results := []map[string]string{}
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return results, err
		}
		curRow := make(map[string]string)
		for k, v := range values {
			curRow[cols[k]] = string(v)
		}
		results = append(results, curRow)
	}
	rows.Close()
	return results, nil
}

func (db *DBServer) Query(queryInfo QueryInfo) []map[string]string {
	if db.DB != nil {
		sql := createQuerySql(queryInfo)
		results, err := dbQueryString(db.DB, sql)
		if err != nil {
			fmt.Println("Query err , sql : ", sql)
		}
		return results
	}
	return nil
}

func (db *DBServer) QuerySQL(sql string) []map[string]string {
	if db.DB != nil {
		results, err := dbQueryString(db.DB, sql)
		if err != nil {
			fmt.Println("Query err : ", err, "sql : ", sql)
		}
		return results
	}
	return nil
}

func (db *DBServer) Exec(execInfo ExecInfo) error {
	if db.DB != nil {
		sqlArr := createExecSql(db.DBtype, execInfo)
		for _, sql := range sqlArr {
			_, err := db.DB.Exec(sql)
			if err != nil {
				fmt.Println("Exec err : ", err, " sql : ", sql)
				return err
			}
			//time.Sleep(10 * time.Microsecond)
		}
	}
	return nil
}

func (db *DBServer) ExecSQL(sql string) error {
	if db.DB != nil {
		_, err := db.DB.Exec(sql)
		if err != nil {
			fmt.Println("Exec err : ", err, " sql : ", sql)
			return err
		}
	}
	return nil
}

func (db *DBServer) TruncateTable(tablename string) error {
	if db.DB != nil {
		found := false
		for _, v := range db.TablesName {
			if v == tablename {
				found = true
			}
		}
		if found {
			sqlstr := "TRUNCATE " + tablename
			_, err := db.DB.Exec(sqlstr)
			checkErr(err)
		}
	}
	return nil
}

func createQuerySql(queryInfo QueryInfo) string {
	sql := ""

	keyStr := ""
	ConditionStr := ""
	// select key1,key2 from tablename where condition1,condition2 order by PK limit 1,100
	if (len(queryInfo.KeyArr) == 0) || (len(queryInfo.KeyArr) == 1 && queryInfo.KeyArr[0] == "*") {
		keyStr = "*"
	} else {
		keyStr = strCombine(queryInfo.KeyArr, ",")
	}
	if len(queryInfo.ConditionArr) > 0 {
		ConditionStr = "where "
		ConditionStr = ConditionStr + strCombine(queryInfo.ConditionArr, " AND ")
	}

	sql = fmt.Sprintf("select %s from %s %s order by %s", keyStr, queryInfo.TableName, ConditionStr, queryInfo.PK)
	return sql
}

func (db *DBServer) GetQueryTotalCount(sql string) uint {
	sqlLower := strings.ToLower(sql)
	pos := strings.Index(sqlLower, "from")
	if pos <= 0 {
		return 0
	}
	totalSql := "select count(*) as count " + sql[pos:]
	orderPos := strings.Index(sqlLower, "order")
	if orderPos > 0 {
		totalSql = "select count(*) as count " + sql[pos:orderPos]
	}
	if totalSql != "" {
		res := db.QuerySQL(totalSql)
		if len(res) > 0 {
			strCount, ok := res[0]["count"]
			if ok == false {
				// for oracle handle.
				strCount = res[0]["COUNT"]
			}
			count, _ := strconv.Atoi(strCount)
			return uint(count)
		}
	}
	return 0
}

func (db *DBServer) GetMaxValue(field string) string {
	sql := ""
	if db.DBtype == ORACLE {
		sql = fmt.Sprintf(" select %s from %s where rownum <= 1 order by %s desc ", field, db.TablesName, field)
	}
	if db.DBtype == MYSQL || db.DBtype == POSTGRESQL {
		sql = fmt.Sprintf(" select %s from %s order by %s desc limit 1 ", field, db.TablesName, field)
	}
	res := db.QuerySQL(sql)
	if len(res) > 0 {
		strCount, ok := res[0][field]
		if ok == false {
			fmt.Println("GetMaxValue error , field : ", field, " not in ", db.TablesName)
		}
		return strCount
	}
	return ""
}

func (db *DBServer) QueryPage(sql string, start int, num int) string {
	if db.DBtype == ORACLE {
		tmpFmt := `
		SELECT *
		FROM (SELECT a.*, ROWNUM rn
				FROM (%s) a
				WHERE ROWNUM <= %d)
		WHERE rn > %d
		`
		return fmt.Sprintf(tmpFmt, sql, (num + start), start)
	}
	if db.DBtype == MYSQL {
		return sql + fmt.Sprintf(" limit %d,%d ", start, num)
	}
	if db.DBtype == POSTGRESQL {
		return sql + fmt.Sprintf(" limit %d offset %d ", num, start)
	}
	return sql
}

func (db *DBServer) GetTableColumns(tableName string) []map[string]string {
	sql := ""
	if db.DBtype == ORACLE {
		strFmt := "select COLUMN_NAME as NAME,DATA_TYPE as TYPE from user_tab_columns where Table_Name='%s' "
		sql = fmt.Sprintf(strFmt, tableName)
	}
	if db.DBtype == POSTGRESQL {
		strFmt := ` SELECT a.attname AS name,t.typname AS type FROM pg_class c,pg_attribute a,pg_type t
		WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid
		ORDER BY a.attnum `
		sql = fmt.Sprintf(strFmt, tableName)
	}
	results, err := dbQueryString(db.DB, sql)
	if err != nil {
		fmt.Println("Query err : ", err, " sql : ", sql)
	}
	return results
}

func createQueryTotalCountSql(queryInfo QueryInfo) string {
	sql := ""
	keyStr := ""
	ConditionStr := ""
	keyStr = " count(*) as count "
	if len(queryInfo.ConditionArr) > 0 {
		ConditionStr = "where "
		ConditionStr = ConditionStr + strCombine(queryInfo.ConditionArr, " AND ")
	}
	sql = fmt.Sprintf("select %s from %s %s order by %s", keyStr, queryInfo.TableName, ConditionStr, queryInfo.PK)
	return sql
}

func createExecSql(dbtype string, execInfo ExecInfo) []string {
	sqlArr := []string{}
	if execInfo.Handle == INSERT {
		sqlArr = insertSql(dbtype, execInfo)
	}
	if execInfo.Handle == UPDATE {
		sqlArr = updateSql(dbtype, execInfo)
	}
	return sqlArr
}

func insertSql(dbtype string, execInfo ExecInfo) []string {
	sql := ""
	sqlArr := []string{}

	results := execInfo.Content
	tablename := execInfo.TableName
	ColumnArr := []string{}
	keyStr, valueStr := "", ""
	for k, _ := range results[0] {
		keyStr += k + ","
		ColumnArr = append(ColumnArr, k)
	}
	keyStr = keyStr[:len(keyStr)-1]

	if dbtype == MYSQL || dbtype == POSTGRESQL || dbtype == SQLITE3 {
		for i := 0; i < len(results); i += MAXSQLCOUNT {
			valueStr = "VALUES"
			endflag := false
			for count := 0; count < MAXSQLCOUNT && !endflag; count++ {
				RowTmpValueStr := "("
				for j := 0; j < len(ColumnArr); j++ {
					if (i + count) >= len(results) {
						endflag = true
						break
					}
					TmpValueStr := results[i+count][ColumnArr[j]]
					if TmpValueStr == "<nil>" {
						TmpValueStr = ""
					}
					if strings.Contains(TmpValueStr, "'") {
						TmpValueStr = strings.Replace(TmpValueStr, "'", "''", -1)
					}
					RowTmpValueStr += "'" + TmpValueStr + "'" + ","
				}
				if endflag {
					RowTmpValueStr = ""
				} else {
					RowTmpValueStr = RowTmpValueStr[:len(RowTmpValueStr)-1]
					RowTmpValueStr += "),"
					valueStr += RowTmpValueStr
				}
			}
			valueStr = valueStr[:len(valueStr)-1]
			sqlstr := "INSERT INTO " + tablename + " (" + keyStr + " ) " + valueStr
			sqlArr = append(sqlArr, sqlstr)
		}
		return sqlArr
	}
	if dbtype == ORACLE {
		beginStr := "BEGIN \r\n"
		endStr := "\r\n END;"

		insertAllStr := ""
		insertStr := fmt.Sprintf(" INSERT INTO %s (%s) ", tablename, keyStr)
		for i := 0; i < len(results); i += MAXSQLCOUNT {
			valueStrFmt := " VALUES (%s) ; "
			valueStr := ""
			endflag := false
			for count := 0; count < MAXSQLCOUNT && !endflag; count++ {
				RowTmpValueStr := ""
				for j := 0; j < len(ColumnArr); j++ {
					if (i + count) >= len(results) {
						endflag = true
						break
					}
					TmpValueStr := results[i+count][ColumnArr[j]]
					if strings.Contains(TmpValueStr, "'") {
						TmpValueStr = strings.Replace(TmpValueStr, "'", "''", -1)
					}
					RowTmpValueStr += "'" + TmpValueStr + "'" + ","
				}
				if endflag {
					RowTmpValueStr = ""
				} else {
					RowTmpValueStr = RowTmpValueStr[:len(RowTmpValueStr)-1]
					valueStr = fmt.Sprintf(valueStrFmt, RowTmpValueStr)
					insertAllStr = insertAllStr + insertStr + valueStr
				}
			}
		}
		sql = beginStr + insertAllStr + endStr
		sqlArr = append(sqlArr, sql)
		return sqlArr
	}
	return sqlArr
}

func updateSql(dbtype string, execInfo ExecInfo) []string {

	// update table set key=value where id=1
	sqlFmt := " update %s set %s where %s='%s' ; "
	content := execInfo.Content
	sqlArr := []string{}
	count := len(execInfo.Content)

	for i := 0; i < count; i += MAXSQLCOUNT {
		sqlTotal := ""
		sql := ""
		endflag := false
		for j := 0; j < MAXSQLCOUNT && !endflag; j++ {
			strKey := ""
			PKValue := ""
			if (i + j) >= count {
				endflag = true
				break
			}
			for k, v := range content[i+j] {
				if k != execInfo.PK {
					strKey = strKey + k + "='" + v + "',"
				} else {
					PKValue = v
				}
			}
			strKey = strKey[:len(strKey)-1]
			sql = fmt.Sprintf(sqlFmt, execInfo.TableName, strKey, execInfo.PK, PKValue)
			sqlTotal = sqlTotal + sql
			sqlArr = append(sqlArr, sqlTotal)
		}
	}
	return sqlArr
}

func strCombine(arr []string, comb string) string {
	res := ""
	for k, v := range arr {
		if k == len(arr)-1 {
			res += v
		} else {
			res += v + comb
		}
	}
	return res
}
