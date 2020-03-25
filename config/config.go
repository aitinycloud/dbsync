//==================================
//  * Name：DataSync
//  * DateTime：2019/08/15
//  * File: config.go
//  * Note: read config file .
//==================================

package config

import (
	"fmt"
	"io/ioutil"
)

var ConfigStringJSON string

type GlobalConfig struct {
	Name  string
	SrcDB map[string]string
	DesDB map[string]string
	Redis map[string]string
}

// 没有main文件，打开的路径也不对
func ParseConfig(configFile string) {
	// Read config JSON.
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("JSON config ReadFile error : %s", err))
	}
	ConfigStringJSON = string(buf)
	fmt.Println(ConfigStringJSON)
}
