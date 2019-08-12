//==================================
//  * Name：DataSync
//  * DateTime：2019/07/22 22:30
//  * File: config.go
//  * Note: read config file .
//==================================

package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/viper"
)

var Config *GlobalConfig

type GlobalConfig struct {
	Name  string
	SrcDB map[string]string
	DesDB map[string]string
	Redis map[string]string
}

// for test.
var ConfigStringJSON = `
{
    "DataSync": [
        {
            "job": {
                "name": "oracle2postgres",
                "srcSql": " SELECT \"id\",\"create_uid\",\"write_date\",\"active\",\"name\",\"price_unit\",\"amount_total\",\"payment_id\",\"create_date\",\"exe_state\" FROM \"his_outpatient_fee\" where  \"id\" > 2986180 and \"id\" < 2986190 and \"create_date\" > to_date('2019-02-04 00:00:00','yyyy-mm-dd hh24:mi:ss') and rownum < 10 ",
                "srcTable": "his_outpatient_fee",
                "desTable": "his_outpatient_fee_new",
                "FieldsMap": {
                    "create_uid": "create_uid_rename",
                    "active": "active_rename",
                    "amount_total": "amount_total_rename"
                },
                "desTablePK": "id",
                "updateType": "incr"
            }
        }
    ]
}
`

// 没有main文件，打开的路径也不对
func ParseConfig(projectName string, configFile string, env string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("ParseConfigError:", err)
		panic(err)
	}
	Config = new(GlobalConfig)
	name := ""
	err = viper.UnmarshalKey("DataSync.Name", &name)
	if err != nil {
		log.Println("UnmarshalConfigError:", err)
		panic(err)
	}

	err = viper.UnmarshalKey("DataSync", &Config)
	if err != nil {
		log.Println("UnmarshalConfigError:", err)
		panic(err)
	}

	// Read config JSON.
	buf, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(fmt.Sprintf("JSON config ReadFile error : ", err))
	}
	ConfigStringJSON = string(buf)
	fmt.Sprintf(ConfigStringJSON)
}
