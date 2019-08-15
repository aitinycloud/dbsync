//==================================

//  * Name：DataSync

//  * DateTime：2019/08/16

//  * File: handle.go

//  * Note: Business processing.

//==================================

package handle

import "testing"

func Test_sqlAddCondition(t *testing.T) {
	type args struct {
		tableName string
		sql       string
		condition string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{"tablename", "select * from tablename", "a=1"}, "select * from tablename where a=1"},
		{"test2", args{"tablename", "select * from tablename where b>=2", "a=1"}, "select * from tablename where a=1 and  b>=2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sqlAddCondition(tt.args.tableName, tt.args.sql, tt.args.condition); got != tt.want {
				t.Errorf("sqlAddCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
