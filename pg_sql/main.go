package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	host     = "x.x.x.x"
	port     = 30822
	user     = "postgres"
	password = "xxxx"
	dbname   = "xxxx"
)

var wg sync.WaitGroup

func writeToLogFile(dbURL, logFileName string) {
	defer wg.Done()

	// 创建或打开日志文件
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 使用log包写入日志到文件
	logger := log.New(file, "[INFO] ", log.LstdFlags)

	var db *sql.DB

	maxAttempts := 5
	for attempts := 0; attempts < maxAttempts; attempts++ {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			logger.Printf("连接数据库时发生错误: %v\n", err)
			logErrorToMySQL(err.Error())
			time.Sleep(time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			logger.Printf("无法连接到数据库: %v\n", err)
			logErrorToMySQL(err.Error())
			time.Sleep(time.Second)
			continue
		}

		logger.Println("pgsql数据库连接成功")
		break
	}

	if db == nil {
		logger.Fatal("无法建立与数据库的连接")
	}
	defer db.Close()

	ticker := time.Tick(30 * time.Second)

	for range ticker {
		sqlline := "SELECT * FROM \"public\".\"l_timelog\" LIMIT 1;"
		rows, err := db.Query(sqlline)
		if err != nil {
			logger.Println("查询时发生错误:", err)
			logErrorToMySQL(err.Error())
			continue
		}

		for rows.Next() {
			var id string
			var className string
			var method string
			var time int
			var CreateDate string
			var message string

			if err := rows.Scan(&id, &className, &method, &time, &CreateDate, &message); err != nil {
				logger.Println("扫描结果时发生错误:", err)
				logErrorToMySQL(err.Error())
				continue
			}
			logger.Printf("id: %s, className: %s\n , method: %s\n ,time: %d\n,CreateDate: %s\n,message: %s\n", id, className, method, time, CreateDate, message)

			/* db2, err := sql.Open("mysql", "root:123456@tcp(192.168.111.137:3306)/info_mysql")
			if err != nil {
				fmt.Println("连接日志数据库时发生错误:", err)
				continue
			}

			// 插入值到第二个数据库的表
			insertSQL := "INSERT INTO mysql_success (name, size) VALUES (?, ?)"
			_, err = db2.Exec(insertSQL, dbRole, dbOrders)
			if err != nil {
				fmt.Println("插入数据时发生错误:", err)
				continue
			}
			fmt.Println("数据成功插入日志数据库的表中") */
		}

		// 连接到第二个数据库

		logger.Println("--------------------------------------------")

		if err := rows.Err(); err != nil {
			logger.Println("遍历结果时发生错误:", err)
			logErrorToMySQL(err.Error())
		}
		rows.Close()
	}
}

// 插入error
func logErrorToMySQL(errorMessage string) {
	var db3 *sql.DB
	db3, err := sql.Open("mysql", "root:123456@tcp(192.168.111.137:3306)/info_pgsql")
	if err != nil {
		log.Println("连接日志数据库时发生错误:", err)
		return
	}

	insertSQL := "INSERT INTO pg_error (name) VALUES (?)"
	_, err = db3.Exec(insertSQL, errorMessage)
	if err != nil {
		log.Println("插入错误日志时发生错误:", err)
	}
	defer db3.Close()
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	wg.Add(2)
	go writeToLogFile(psqlInfo, "PGsql38200.log")
	wg.Wait()
}
