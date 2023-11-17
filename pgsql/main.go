package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	_ "github.com/lib/pq"
)

const (
	host     = "IP地址"
	port     = 端口
	user     = "用户名"
	password = "密码"
	dbname   = "数据库名"
)


func init(){
	//创建日志文件
	logfile, err := os.Create("PgLogs2.log")
	if err != nil {
		log.Fatal("未成功创建文件:", err)
	}
	defer logfile.Close()


	//写入日志文件                                                                                                                                                                                         
    writefile,_ := os.OpenFile("PgLogs2.log",os.O_RDWR | os.O_CREATE | os.O_APPEND,0766)
	log.SetOutput(writefile)
	log.SetPrefix("[INFO]")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	// 连接数据库
	for {
		db, err = sql.Open("postgres", psqlInfo)
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

		log.Println("PostgreSQL数据库连接成功")
		break 
	}

	if db == nil {
		log.Fatal("无法建立与数据库的连接")
	}
	defer db.Close()



	ticker := time.Tick(30 * time.Second)  //每30执行一次
	for range ticker {
		sqlline := "语句;"
		rows, err := db.Query(sqlline)
		if err != nil {
			log.Println("查询时发生错误:", err)
			continue
		}

		for rows.Next() {
			var a string
			var b string
			var c string
			if err := rows.Scan(&id,&b,&c); err != nil {
				log.Println("扫描结果时发生错误:", err)
				continue
			}
			log.Printf("a: %s, b: %s\n , c: %s\n",a,b,c)
		}
		log.Println("--------------------------------------------")

		if err := rows.Err(); err != nil {
			log.Println("遍历结果时发生错误:", err)
		}
	}    
}