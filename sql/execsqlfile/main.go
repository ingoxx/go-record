package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	port     = flag.Int("port", 1433, "The sqlserver port")
	host     = flag.String("host", "", "The sqlserver host")
	user     = flag.String("user", "", "The sqlserver user")
	db       = flag.String("db", "", "The sqlserver db")
	password = flag.String("password", "", "The sqlserver password")
	sqlFile  = flag.String("sqlfile", "", "The sqlserver sqlFile")
	query    = flag.String("query", "", "sql cmd")
)

func main() {
	flag.Parse()

	if flag.NFlag() <= 4 {
		log.Fatalln(flag.ErrHelp)
	}

	encodedPassword := url.QueryEscape(*password)

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable&database=%s&timeout=10", *user, encodedPassword, *host, *port, *db)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database >>>", err)
	}
	defer db.Close()

	if *query != "" && *sqlFile != "" {
		log.Fatalln("choose one of param between --sqlfile and --query")
	}

	if *query != "" {
		exec, err := db.Exec(*query)
		if err != nil {
			log.Fatalln("Failed to execute SQL >>>", err)
		}
		rs, err := exec.RowsAffected()
		if err != nil {
			log.Fatalln("Failed to obtain execution results >>>", err)
		}
		log.Printf("query res %d", rs)
	}

	if *sqlFile != "" {
		fo, err := os.Stat(*sqlFile)
		if err != nil {
			log.Fatalln("Failed to open SQL file >>>", err)
		}
		if fo.IsDir() {
			log.Fatalln("Expected a SQL file, not a directory")
		}

		content, err := os.ReadFile(*sqlFile)
		if err != nil {
			log.Fatalln("Failed to read SQL file >>>", err)
		}

		batches := splitSQLBatches(string(content))
		for idx, stmt := range batches {
			if strings.TrimSpace(stmt) == "" {
				continue
			}
			_, err := db.Exec(stmt)
			if err != nil {
				log.Fatalf("Error executing batch %d >>> %v\nSQL:\n%s", idx+1, err, stmt)
			}
			log.Printf("Batch %d executed successfully", idx+1)
		}
	}
}

func splitSQLBatches(sqlText string) []string {
	lines := strings.Split(sqlText, "\n")
	var batches []string
	var current []string

	for _, line := range lines {
		if strings.ToUpper(strings.TrimSpace(line)) == "GO" {
			if len(current) > 0 {
				batches = append(batches, strings.Join(current, "\n"))
				current = []string{}
			}
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 {
		batches = append(batches, strings.Join(current, "\n"))
	}

	return batches
}
