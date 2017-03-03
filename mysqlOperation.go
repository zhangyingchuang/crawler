// mysqlOperation
package main

import (
	"database/sql"
	"fmt"

	"io/ioutil"
	"log"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

//数据库中的表以cj_98_ 和cj_99_ 开头 后面为日期
//日期格式为： 2015_05_28
const (
	startName1 = "cj_98_"
	startName2 = "cj_99_"
	backUpPath = "/data/caiji_bak" //备份的路径

	//mysql 配置信息
	user   = "myuser"
	passwd = "mypassword"
	host   = "192.168.0.32"
	db     = "zp_mysql"
	path   = "./"
)

var backups []string //存储已经备份过的表名

var cj_98_no int = 0 //记录最新cj_98的第一张表
var cj_99_no int = 0 //记录最新cj_99的第一张表

func init() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	loadBackup()
}

func loadBackup() {

	// 加载数据表cj_backup_record 表 把 已经备份过的表名写入数组中
	log.Println("把已经备份过的表名写入数组中")

	// 连接数据库
	mysql, err := sql.Open("mysql", loadconfig())
	if err != nil {
		log.Fatalln("数据打开失败 ，请检查数据配置", err)
	}
	defer mysql.Close()

	var tmp string = ""

	rows, err := mysql.Query("select table_name from cj_backup_record")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&tmp); err != nil {
			log.Fatal(err)
		}
		//
		backups = append(backups, tmp)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

//从2015.5.1开始
func dealDatabase() {

	//1 连接数据库
	mysql, err := sql.Open("mysql", loadconfig())
	if err != nil {
		log.Fatalln("数据打开失败 ，请检查数据配置", err)
	}
	defer mysql.Close()

	t := time.Date(2015, time.May, 1, 0, 0, 0, 0, time.UTC)
	var end int64 = t.Unix()

	start := time.Now().Unix()
	for {
		if start <= end {
			fmt.Println("此次备份结束！")
			break
		}
		fmt.Println("检查cj表：", getDate(start))
		name1 := startName1 + getDate(start)

		name2 := startName2 + getDate(start)

		//1. 检查表是否存在 2. 判断表是否已经备份过
		if checkTable(name1, mysql) {

			cj_98_no++
			//log.Println("test:", cj_98_no)
			if cj_98_no != 1 {
				if !checkBachup(name1) {
					fmt.Println(name1, "表存在开始备份")
					backUpTable(name1, mysql)
				}

			} else {
				log.Println(name1, "为cj_98 的第一张表 不做备份")
			}

		}

		if checkTable(name2, mysql) {

			cj_99_no++
			if cj_99_no != 1 {
				if !checkBachup(name2) {
					fmt.Println(name2, "表存在开始备份")
					backUpTable(name2, mysql)
				}

			} else {
				log.Println(name2, "为cj_99 的第一张表 不做备份")
			}

		}

		start -= 24 * 60 * 60
	}

}

//获取日期 如：2015_05_28
func getDate(mytime int64) string {

	date := strings.Split(time.Unix(mytime, 0).String(), " ")[0]
	return strings.Replace(date, "-", "_", -1)
}

//检查表是否存在
func checkTable(name string, mysql *sql.DB) bool {

	sql := "show tables like" + `'%` + name + `%'`
	rows, err := mysql.Query(sql)

	if err != nil {
		debug.PrintStack()
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}

//执行shell 脚本
func excuteShell(cmdstr string) bool {

	cmd := exec.Command("/bin/sh", "-c", cmdstr)
	log.Println("执行命令：", cmdstr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("StdoutPipe: " + err.Error())

	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("StderrPipe: ", err.Error())

	}

	if err := cmd.Start(); err != nil {
		log.Println("Start: ", err.Error())

	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println("ReadAll stderr: ", err.Error())

	}

	if len(bytesErr) != 0 {
		log.Printf("stderr is not nil: %s", bytesErr)

	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll stdout: ", err.Error())

	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: ", err.Error())
		return false

	}
	log.Println(string(bytes))
	return true

}

func backUpTable(name string, mysql *sql.DB) {

	//1.导出文件为.sql
	str :=
		"mysqldump -u " + user + " " + "-p" + passwd + " " +
			"-h" + " " + host + " " + db + " " + name + "" + "> " + path + name + ".sql"

	if !excuteShell(str) {
		log.Fatalln("导出数据库出错!")
	}

	//tar -zcvf $i.tar.gz $i;
	//2.压缩文件在当前目录
	tarstr := "tar -zcvf " + name + ".tar.gz" + " " + name + ".sql"
	if !excuteShell(tarstr) {
		sendEmail("压缩文件出错", "压缩文件出错")
		log.Fatalln("压缩文件出错!")
	}

	//3.mv 到指定目录
	mvstr := "mv " + name + ".tar.gz " + backUpPath
	if !excuteShell(mvstr) {
		sendEmail("移动文件出错", "移动文件出错")
		log.Fatalln("移动文件出错!")
	}

	// 4.删除.sql 文件
	rmstr := "rm -rf " + " " + name + ".sql"
	if !excuteShell(rmstr) {
		sendEmail("删除文件出错", "删除文件出错")
		log.Fatalln("删除文件出错!")
	}

	//5. 删除数据库中的表
	////deleteTable(name, mysql)

	//6. 记录操作信息
	logOperation(mysql, name)
	fmt.Println(name, "备份完成")

}

//删除表
func deleteTable(name string, mysql *sql.DB) {
	sql := "DROP TABLE " + name
	rows, err := mysql.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
}

//记录操作
func logOperation(mysql *sql.DB, name string) {
	count := getTableCount(mysql, name)

	stmt, err := mysql.Prepare("insert into cj_backup_record (date,backup_name,table_name,count) values(now(),?,?,?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(name+".tar.gz", name, count)

	checkErr(err)

}

//查询表的记录条数
func getTableCount(db *sql.DB, name string) int {
	var count = 0
	rows, err := db.Query("select * from " + name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// 读取数据
	for rows.Next() {
		count++
	}

	return count
}
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

//检查cj_backup_record 表 检查表是否已经备份过
func checkBachup(name string) bool {
	for _, value := range backups {
		if name == value {
			return true
		}
	}
	return false
}
