// all_test
package main

import (
	"testing"
)

func Test_loadBackup(t *testing.T) {

	if checkBachup("cj_98_2015_05_01") {
		t.Log("表存在")
	} else {
		t.Log("表不存在")
	}
}
