//==================================
//  * Name：DataSync
//  * DateTime：2019/07/22 22:30
//  * File: handle.go
//  * Note: Business processing.
//==================================

package handle

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"../config"
	"../db"
	"../pkg/cache"
	"../pkg/logging"
	"github.com/go-redis/redis"
	gcache "github.com/patrickmn/go-cache"
)

var SrcDBPtr *db.DBServer
var DesDBPtr *db.DBServer

var redisclient *redis.Client

func Setup() {
	logging.Info(fmt.Sprintf("DataSync handle Setup , read config ."))
	config.ParseConfig("DataSync", "config.yaml", "")
	//cache init.
	cache.CacheInit()
}

func Work() {
	logging.Info(fmt.Sprintf("DataSync handle work."))
	db.Setup()
	//Get config.
	SrcDB := db.DBServer{}
	
	SrcDB.DBtype = gjson.Get(config.ConfigStringJSON, "SrcDB.type").String()
	SrcDB.User = gjson.Get(config.ConfigStringJSON, "SrcDB.User").String()
	SrcDB.Passwd = gjson.Get(config.ConfigStringJSON, "SrcDB.Passwd").String()
	SrcDB.Host = gjson.Get(config.ConfigStringJSON, "SrcDB.Host").String()
	SrcDB.DBName = gjson.Get(config.ConfigStringJSON, "SrcDB.DBName").String()
	err := SrcDB.Start()
	if err != nil {
		logging.Info(fmt.Sprintf("SrcDB Link Fail,Please check SrcDB config."))
		panic("SrcDB Link Fail,Please check SrcDB config.")
	}
	SrcDBPtr = &SrcDB
	DesDB := db.DBServer{}
	DesDB.DBtype = gjson.Get(config.ConfigStringJSON, "DesDB.type").String()
	DesDB.User = gjson.Get(config.ConfigStringJSON, "DesDB.User").String()
	DesDB.Passwd = gjson.Get(config.ConfigStringJSON, "DesDB.Passwd").String()
	DesDB.Host = gjson.Get(config.ConfigStringJSON, "DesDB.Host").String()
	DesDB.DBName = gjson.Get(config.ConfigStringJSON, "DesDB.DBName").String()
	err = DesDB.Start()
	if err != nil {
		logging.Info(fmt.Sprintf("DesDB Link Fail,Please check DesDB config."))
		panic("DesDB Link Fail,Please check DesDB config.")
	}
	DesDBPtr = &DesDB

	jobMap := gjson.Get(config.ConfigStringJSON, "DataSync.0.job").Map()
	sql := jobMap["srcSql"].String()
	//get FieldsMap config.
	reNameConfig := jobMap["FieldsMap"].Map()
	reNameMapInfo := db.ReNameMapInfo{}
	reNameMapInfo.SrcTableName = jobMap["srcTable"].String()
	reNameMapInfo.DesTableName = jobMap["desTable"].String()
	NameMap := make(map[string]string)
	for k, v := range reNameConfig {
		NameMap[k] = v.String()
	}
	reNameMapInfo.NameMap = NameMap
	// get srcDB tableColumns
	tableColumnsResult := SrcDB.GetTableColumns(jobMap["srcTable"].String())
	tableColumnsMap := make(map[string]string)
	for i := 0; i < len(tableColumnsResult); i++ {
		key := tableColumnsResult[i]["NAME"]
		val := tableColumnsResult[i]["TYPE"]
		if res, ok := NameMap[key]; ok {
			tableColumnsMap[res] = val
		} else {
			tableColumnsMap[key] = val
		}
	}
	strColumns, _ := json.Marshal(tableColumnsMap)
	logging.Info(fmt.Sprintf("tableColumnsMap : %s", string(strColumns)))
	//get Des Table to redis.
	srcTableName := jobMap["srcTable"].String()
	desTableName := jobMap["desTable"].String()
	sql = "select count(*) as count from " + desTableName
	desTotalCount := DesDB.GetQueryTotalCount(sql)

	for i := 0; i < int(desTotalCount); i += db.MAXPAGECOUNT {
		sql := "select * from " + desTableName
		sql = DesDBPtr.QueryPage(sql, i, db.MAXPAGECOUNT)
		desTableResult := DesDBPtr.QuerySQL(sql)
		CacheToLocal(desTableName, desTableResult)
	}

	//support page query. get total count .
	sql = jobMap["srcSql"].String()
	totalCount := SrcDB.GetQueryTotalCount(sql)
	logging.Info(fmt.Sprintf("Source tableName : %s , totalCount : %d ", srcTableName, totalCount))
	logging.Info(fmt.Sprintf("Destination tableName : %s , totalCount : %d", desTableName, desTotalCount))
	if totalCount > db.MAXPAGECOUNT {
		start := 0
		num := db.MAXPAGECOUNT
		for start = 0; start < int(totalCount); start += db.MAXPAGECOUNT {
			logging.Info(fmt.Sprintf("Paging query . start : %d,num : %d. handled : %d .", start, num, start))
			sqlQueryPage := SrcDB.QueryPage(sql, start, num)
			//
			srcTableResult := SrcDB.QuerySQL(sqlQueryPage)
			renameResult := SrcDB.ReName(srcTableResult, reNameMapInfo)
			//renameResultstr, _ := json.Marshal(renameResult)
			//logging.Info("ReNameResult : ", string(renameResultstr))
			insertExec, updateExec := CompareWithCache(tableColumnsMap, renameResult)
			InsertAndUpdate(insertExec, updateExec)
		}
	} else {
		//
		srcTableResult := SrcDB.QuerySQL(sql)
		renameResult := SrcDB.ReName(srcTableResult, reNameMapInfo)
		//renameResultstr, _ := json.Marshal(renameResult)
		//logging.Info("ReNameResult : ", string(renameResultstr))
		insertExec, updateExec := CompareWithCache(tableColumnsMap, renameResult)
		InsertAndUpdate(insertExec, updateExec)
	}
}

func CacheToLocal(tableName string, desTableResult []map[string]string) {
	jobMap := gjson.Get(config.ConfigStringJSON, "DataSync.0.job").Map()
	tablePK := jobMap["desTablePK"].String()
	for _, v := range desTableResult {
		key := v[tablePK]
		strVal, _ := json.Marshal(v)
		cache.GCache.Set(tableName + "_" + key, string(strVal), gcache.NoExpiration)
	}
}

func CompareWithCache(tableColumnsMap map[string]string, renameResult []map[string]string) (db.ExecInfo, db.ExecInfo) {
	// Compare.
	jobMap := gjson.Get(config.ConfigStringJSON, "DataSync.0.job").Map()
	tableName := jobMap["desTable"].String()
	tablePK := jobMap["desTablePK"].String()

	keysArr := []string{}
	for k, _ := range renameResult[0] {
		keysArr = append(keysArr, k)
	}
	logging.Info(keysArr)
	insertExec := db.ExecInfo{DesDBPtr.DBName, tableName, db.INSERT, tablePK, []map[string]string{}}
	updateExec := db.ExecInfo{DesDBPtr.DBName, tableName, db.UPDATE, tablePK, []map[string]string{}}
	for _, v := range renameResult {
		handle := ""
		srcQueryStr := ""
		cacheQueryStr := ""
		srcQueryMD5 := ""
		cacheQueryMD5 := ""

		for i := 0; i < len(keysArr); i++ {
			if !isCompare(tableColumnsMap[keysArr[i]]) {
				continue
			}
			str := v[keysArr[i]] + ";"
			if res, resstr := TypeConv(tableColumnsMap[keysArr[i]], v[keysArr[i]]); res {
				str = resstr + ";"
			}
			srcQueryStr = srcQueryStr + str
			valJSONStr, _ := redisclient.HGet(tableName, v[tablePK]).Result()
			valMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(valJSONStr), &valMap)
			if err != nil {
				// key is not exist. Inset hande.
				handle = db.INSERT
				break
			}
			val := valMap[keysArr[i]].(string)
			str = val + ";"
			cacheQueryStr = cacheQueryStr + str
		}
		strShow, _ := json.Marshal(v)
		logging.Info(fmt.Sprintf("need handle : %s , value :  %s, srcQueryStr : %s , cacheQueryStr : %s", handle, strShow, srcQueryStr, cacheQueryStr))
		if handle == "" {
			srcQueryMD5 = fmt.Sprintf("%x", md5.Sum([]byte(srcQueryStr)))
			cacheQueryMD5 = fmt.Sprintf("%x", md5.Sum([]byte(cacheQueryStr)))
			logging.Info(fmt.Sprintf("srcQueryMD5 : %s , cacheQueryMD5 : %s", srcQueryMD5, cacheQueryMD5))
			if srcQueryMD5 != cacheQueryMD5 {
				// update handle.
				handle = db.UPDATE
			}
		}

		if handle == db.INSERT {
			insertExec.Content = append(insertExec.Content, v)
		}
		if handle == db.UPDATE {
			updateExec.Content = append(updateExec.Content, v)
		}
	}
	return insertExec, updateExec
}

func InsertAndUpdate(insertExec db.ExecInfo, updateExec db.ExecInfo) {
	//insert handle
	//insertStr, _ := json.Marshal(insertExec)
	//updateStr, _ := json.Marshal(updateExec)
	//logging.Info(fmt.Sprintf("insertExec : %s\nupdateExec:%s\n", string(insertStr), string(updateStr)))

	if len(insertExec.Content) > 0 {
		logging.Info(fmt.Sprintf("insertExec Count : %d ", len(insertExec.Content)))
		err := DesDBPtr.Exec(insertExec)
		if err != nil {
			logging.Info(err)
		}
		//update redis.
		CacheToLocal(insertExec.TableName, insertExec.Content)
	}
	//update handle.
	if len(updateExec.Content) > 0 {
		logging.Info(fmt.Sprintf("updateExec Count : %d ", len(updateExec.Content)))
		err := DesDBPtr.Exec(updateExec)
		if err != nil {
			logging.Info(err)
		}
		//update redis.
		CacheToLocal(insertExec.TableName, insertExec.Content)
	}
}

func isCompare(nameType string) bool {
	if nameType == "DATE" || nameType == "timestamp" {
		return false
	}
	return true
}

func TypeConv(nameType string, value string) (bool, string) {
	valueStr := strings.ToLower(value)
	//bool
	if (nameType == "VARCHAR2") && (valueStr == "t" || valueStr == "f") {
		if valueStr == "t" {
			return true, "true"
		}
		if valueStr == "f" {
			return true, "false"
		}
	}
	//DATE
	if nameType == "DATE" {

	}
	return false, ""
}
