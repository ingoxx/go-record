package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	port     = flag.Int("port", 1433, "The sqlserver port")
	host     = flag.String("host", "", "The sqlserver host")
	user     = flag.String("user", "", "The sqlserver user")
	password = flag.String("password", "", "The sqlserver password")
	sqlFile  = flag.String("sqlfile", "", "The sqlserver sqlFile")
	query    = flag.String("query", "", "sql cmd")
)

func main() {

	flag.Parse()

	if flag.NFlag() <= 4 {
		log.Fatalln(flag.ErrHelp)
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable&timeout=5", *user, *password, *host, *port)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database >>>", err)
	}

	defer db.Close()

	if *query != "" && *sqlFile != "" {
		log.Fatalln("choose one of param between --sqlFile,--query")
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
		if err == nil {
			if !fo.IsDir() {
				b, err := os.ReadFile(*sqlFile)
				if err != nil {
					log.Fatalln("Can only execute SQL files >>>", err)
				}

				r, err := db.Exec(string(b))
				if err != nil {
					log.Fatalln("Failed to execute SQL >>>", err)
				}

				rs, err := r.RowsAffected()
				if err != nil {
					log.Fatalln("Failed to obtain execution results >>>", err)
				}

				log.Printf("res %d", rs)
			}
		} else {
			log.Fatalln("Can only execute SQL files >>>", err)
		}
	}

}
