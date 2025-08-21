package main

import (
	"database/sql"
	"fmt"
	"log"
)
import _ "github.com/go-sql-driver/mysql"

func checkClusterExists(db *sql.DB) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM cluster_models WHERE master_ip = ?)"
	err := db.QueryRow(query, "2.2.2.2").Scan(&exists)
	if err != nil {
		log.Printf("Failed to check cluster existence: %v", err)
		return exists
	}

	return exists
}

func main() {
	db, err := sql.Open("mysql", "root:7109667@Lxb@tcp(127.0.0.1:34306)/goweb?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	res := checkClusterExists(db)
	fmt.Println("res >>> ", res)
}
