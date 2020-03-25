//==================================
//  * Name：DataSync
//  * DateTime：2019/07/22 22:30
//  * File: main.go
//  * Note: main handle .
//==================================

package main

import (
	"dbsync/handle"
)

func main() {
	handle.Setup()
	handle.Work()
}
