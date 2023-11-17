package main

import (
	"database/sql"
	_ "fmt"
	"log"
	"os"
	"time"
	_ "github.com/go-sql-driver/mysql"
)



func init(){
	//创建日志文件
	logfile, err := os.Create("MySQL2.log")
	if err != nil {
		log.Fatal("未成功创建文件:", err)
	}
	defer logfile.Close()


	//写入日志文件                                                                                                                                                                                         
    writefile,_ := os.OpenFile("MySQL2.log",os.O_RDWR | os.O_CREATE | os.O_APPEND,0766)
	log.SetOutput(writefile)
	log.SetPrefix("[INFO]")
}

func main() {

	var db *sql.DB
	var err error

	// 连接数据库
	for {
		db, err = sql.Open("mysql", "账号:密码@tcp(ip:端口)/数据库")
		if err != nil {
			log.Printf("连接数据库时发生错误 :%v\n", err)
			time.Sleep(time.Second) 
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("无法连接到数据库 : %v\n", err)
			time.Sleep(time.Second) 
			continue
		}

		log.Println("Mysql数据库连接成功")
		break 
	}

	if db == nil {
		log.Fatal("无法建立与数据库的连接")
	}
	defer db.Close()



	ticker := time.Tick(30 * time.Second)
	for range ticker {
		sqlline := "语句;"
		rows, err := db.Query(sqlline)
		if err != nil {
			log.Println("查询时发生错误:", err)
			continue
		}

		for rows.Next() {
			var a int
			var b string
			var c string
			if err := rows.Scan(&a ,&b, &c); err != nil {
				log.Println("扫描结果时发生错误:", err)
				continue
			}
			log.Printf("a : %d b: %s, c: %s\n",a, b, c)
		}
		log.Println("--------------------------------------------")

		if err := rows.Err(); err != nil {
			log.Println("遍历结果时发生错误:", err)
		}
	}    
}