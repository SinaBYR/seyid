package lib

import "database/sql"

func InitDatabase () *sql.DB {
	connStr := "user=postgres dbname=seyid password=Sina13801111 host=localhost port=5432 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

