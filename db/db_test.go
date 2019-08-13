package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
)

func TestStrCombine(t *testing.T) {
	type args struct {
		arr  []string
		comb string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"t1", args{[]string{}, ""}, ""},
		{"t2", args{[]string{"key1", "key2"}, ","}, "key1,key2"},
		{"t3", args{[]string{"cond1", "cond2"}, " AND "}, "cond1 AND cond2"},
		{"t4", args{[]string{"1", "2"}, ","}, "1,2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strCombine(tt.args.arr, tt.args.comb); got != tt.want {
				t.Errorf("strCombine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateQuerySql(t *testing.T) {
	type args struct {
		queryInfo QueryInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{queryInfo: QueryInfo{}}, want: ""},
		{name: "test2", args: args{queryInfo: QueryInfo{}}, want: ""},
	}
	//test1
	tests[0].args.queryInfo.DBName = "testdb"
	tests[0].args.queryInfo.TableName = "testtable"
	tests[0].args.queryInfo.PK = "id"
	KeyArr := []string{"*"}
	tests[0].args.queryInfo.KeyArr = KeyArr
	tests[0].want = "select * from testtable  order by id"

	//test2
	tests[1].args.queryInfo.DBName = "testdb"
	tests[1].args.queryInfo.TableName = "testtable"
	tests[1].args.queryInfo.PK = "id"
	KeyArr = []string{"a", "b", "c", "d"}
	tests[1].args.queryInfo.KeyArr = KeyArr
	ConditionArr := []string{"a=1", "b<4", "c=d", "d>1"}
	tests[1].args.queryInfo.ConditionArr = ConditionArr
	tests[1].want = "select * from testtable  order by id"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createQuerySql(tt.args.queryInfo); got != tt.want {
				t.Errorf("createQuerySql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBServerQuery(t *testing.T) {

	/*
		oracleServer := DBServer{}
		oracleServer.DBtype = "oracle"
		oracleServer.Host = "192.168.0.80:1521"
		oracleServer.User = "scott"
		oracleServer.Passwd = "scott"
		oracleServer.DBName = "SCOTT"
		err := oracleServer.Start()
		if err != nil {
			t.Errorf("%s", err)
		}
		queryInfo := QueryInfo{}
		queryInfo.DBName = "SCOTT"
		queryInfo.TableName = "DEPT"
		queryInfo.PK = "DEPTNO"
		//KeyArr := []string{"DEPTNO", "DNAME", "LOC"}
		//queryInfo.KeyArr = KeyArr
		//ConditionArr := []string{"DEPTNO>10", "LOC='BOSTON'"}
		//queryInfo.ConditionArr = ConditionArr

		results := oracleServer.Query(queryInfo)
		if results == nil {
			t.Errorf("Query Fail.")
		}
	*/
	postgresServer := DBServer{}
	postgresServer.DBtype = "postgres"
	postgresServer.Host = "192.168.0.80:5432"
	postgresServer.User = "postgres"
	postgresServer.Passwd = "postgres"
	postgresServer.DBName = "testdb"
	err := postgresServer.Start()
	if err != nil {
		t.Errorf("%s", err)
	}
	queryInfo := QueryInfo{}
	queryInfo.DBName = "testdb"
	queryInfo.TableName = "his_clinic_item_category"
	queryInfo.PK = "id"
	KeyArr := []string{"id", "create_uid", "code", "name", "write_date", "his_id"}
	queryInfo.KeyArr = KeyArr
	ConditionArr := []string{"create_date > '2019-04-27'"}
	queryInfo.ConditionArr = ConditionArr

	results := postgresServer.Query(queryInfo)
	if results == nil {
		t.Errorf("Query Fail.")
	}
	t.Log(results)

}

func TestDBServerReName(t *testing.T) {
	type args struct {
		resultMap  []map[string]string
		reNameInfo ReNameMapInfo
	}
	type TestCase struct {
		name string
		db   *DBServer
		args args
		want []map[string]string
	}
	tests := []TestCase{}
	testItem := TestCase{}
	testItem.name = "test"
	testItem.db = nil
	testItemArgs := args{}
	resultMap := []map[string]string{
		{"id": "1", "name": "xiaoming", "age": "22", "address": "chongqing"},
		{"id": "2", "name": "xiaohong", "age": "24", "address": "beijing"},
	}
	testItemArgs.resultMap = resultMap
	reNameInfo := ReNameMapInfo{}
	reNameInfo.SrcTableName = "srvTable"
	reNameInfo.DesTableName = "desTable"
	nameMap := make(map[string]string)
	nameMap["age"] = "curAge"
	nameMap["address"] = "Address"
	reNameInfo.NameMap = nameMap

	testItemArgs.reNameInfo = reNameInfo
	testItem.args = testItemArgs
	testItem.want = nil
	tests = append(tests, testItem)

	t.Log(testItemArgs.resultMap)
	got := testItem.db.ReName(testItem.args.resultMap, testItem.args.reNameInfo)
	t.Log(got)
	if got == nil {
		t.Errorf("DBServer.ReName() Error")
	}
}

func TestDBServerExec(t *testing.T) {
	type args struct {
		execInfo ExecInfo
	}
	type TestCase struct {
		name    string
		db      *DBServer
		args    args
		wantErr bool
	}
	tests := []TestCase{
		// TODO: Add test cases.
	}

	testcase := TestCase{}
	testcase.name = "test1"
	testcase.db = nil
	_ = tests
	/*
		targs := args{}
		execInfo := ExecInfo{}
		execInfo.DBName = "testdb"
		execInfo.TableName = "testtable"
		execInfo.Handle = INSERT
		content := []map[string]string{}
		for i := 0; i < 5; i++ {
			itemMap := make(map[string]string)
			itemMap["ID"] = fmt.Sprintf("%d", i+10)
			itemMap["NAME"] = fmt.Sprintf("name_%d", i)
			content = append(content, itemMap)
		}
		execInfo.Content = content

		targs.execInfo = execInfo
		testcase.args = targs
		testcase.wantErr = false
		_ = tests

		postgresServer := DBServer{}
		postgresServer.DBtype = "postgres"
		postgresServer.Host = "192.168.0.80:5432"
		postgresServer.User = "postgres"
		postgresServer.Passwd = "postgres"
		postgresServer.DBName = "testdb"
		err := postgresServer.Start()
		if err != nil {
			t.Errorf("%s", err)
		}
		t.Log(execInfo)
		postgresServer.Exec(execInfo)
	*/

	oracleServer := DBServer{}
	oracleServer.DBtype = "oracle"
	oracleServer.Host = "192.168.0.80:1521"
	oracleServer.User = "scott"
	oracleServer.Passwd = "scott"
	oracleServer.DBName = "SCOTT"
	err := oracleServer.Start()
	if err != nil {
		t.Errorf("%s", err)
	}
	execInfo := ExecInfo{}
	execInfo.DBName = "SCOTT"
	execInfo.TableName = "TEST3"
	execInfo.Handle = INSERT
	content := []map[string]string{}
	for i := 0; i < 5; i++ {
		itemMap := make(map[string]string)
		itemMap["ID"] = fmt.Sprintf("%d", i+20)
		itemMap["NAME"] = fmt.Sprintf("name_%d", i)
		content = append(content, itemMap)
	}
	execInfo.Content = content
	oracleServer.Exec(execInfo)
}

func TestDBServerExecReNameHandle(t *testing.T) {
	postgresServer := DBServer{}
	postgresServer.DBtype = "postgres"
	postgresServer.Host = "192.168.0.80:5432"
	postgresServer.User = "postgres"
	postgresServer.Passwd = "postgres"
	postgresServer.DBName = "testdb"
	err := postgresServer.Start()
	if err != nil {
		t.Errorf("%s", err)
	}
	queryInfo := QueryInfo{}
	queryInfo.DBName = "testdb"
	queryInfo.TableName = "his_outpatient_fee"
	queryInfo.PK = "id"
	KeyArr := []string{"id", "create_uid", "write_date", "active", "name", "price_unit", "amount_total", "payment_id"}
	queryInfo.KeyArr = KeyArr
	ConditionArr := []string{"id > 2986180", "id < 2986190", "create_date > '2019-05-1'"}
	queryInfo.ConditionArr = ConditionArr

	results := postgresServer.Query(queryInfo)
	if results == nil {
		t.Errorf("Query Fail.")
	}
	t.Log(results)
	reNameInfo := ReNameMapInfo{}
	reNameInfo.SrcTableName = "his_outpatient_fee"
	reNameInfo.DesTableName = "his_outpatient_fee_new"
	nameMap := make(map[string]string)
	nameMap["create_uid"] = "create_uid_rename"
	nameMap["active"] = "active_rename"
	nameMap["amount_total"] = "amount_total_rename"
	reNameInfo.NameMap = nameMap
	resMap := postgresServer.ReName(results, reNameInfo)
	t.Log(resMap)

	execInfo := ExecInfo{}
	execInfo.DBName = "testdb"
	execInfo.TableName = "his_outpatient_fee_new"
	execInfo.Handle = INSERT
	execInfo.Content = resMap
	err = postgresServer.Exec(execInfo)
	if err != nil {
		t.Errorf("Exec fail.")
	}
}

func TestUpdateSql(t *testing.T) {
	dbtype := POSTGRESQL
	execInfo := ExecInfo{}
	execInfo.DBName = "insight"
	execInfo.TableName = "DestinationNode"
	execInfo.PK = "id"
	execInfo.Handle = UPDATE
	content := []map[string]string{}
	resultMap := make(map[string]string)
	resultMap["id"] = "8"
	resultMap["shortName"] = "sName8"
	resultMap["type"] = "8"
	content = append(content, resultMap)

	resultMap2 := make(map[string]string)
	resultMap2["id"] = "9"
	resultMap2["shortName"] = "sName9"
	resultMap2["type"] = "9"
	content = append(content, resultMap2)

	execInfo.Content = content
	res := updateSql(dbtype, execInfo)
	fmt.Println(res)
}
