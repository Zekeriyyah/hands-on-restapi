package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbdriver *sql.DB) {
	//Initialize table train
	statement, driverError := dbdriver.Prepare(train)
	if driverError != nil {
		log.Println(driverError)
	}
	statement.Exec()

	//Initialize table station
	statement, stmtErr := dbdriver.Prepare(station)
	if stmtErr != nil {
		log.Println(stmtErr)
	}
	statement.Exec()

	//Initialize table schedule
	statement, stmtErr = dbdriver.Prepare(schedule)
	if stmtErr != nil {
		log.Println(stmtErr)
	}
	statement.Exec()

	log.Println("Successfully Created/Initialized all tables!!")
}
