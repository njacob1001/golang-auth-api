package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {


	//postgresURI := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", "rummomega", authenticationToken, "rumm-labs.chky1xowzqwg.us-east-2.rds.amazonaws.com", "5432", "rumm-labs")
	postgresURI := fmt.Sprintf("user=%v dbname=%v password=%v host=%v sslmode=disable", "rummomega", "rummlabs","yM0j+WzI9R;bvZmlx^TjsVl,}", "rumm-labs.chky1xowzqwg.us-east-2.rds.amazonaws.com")


	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

