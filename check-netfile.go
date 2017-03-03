// check-netfile
package main

import (
	"os"
)

var checkPath = "/data/caiji_bak/test.txt"

func checkNetFile() bool {
	return Exist(checkPath)
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
