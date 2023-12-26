package go_utils

import (
	"log"
	"database/sql"

	_ "github.com/lib/pq"
)

type PGClient interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// 	connectionString -> "user=postgres password=postgres dbname=postgres sslmode=disable"
func BuildPostgresClient(connectionString string) (PGClient, error) {

	db, err := sql.Open("postgres", connectionString)
	log.Println("ERROR opening postgres connection with github.com/neurocollective/go_utils.BuildPostgresClient() ->")
	log.Println(err.Error())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// takes a struct-specific `scanRows`
func QueryForStructs[T any](
	client PGClient, 
	scanRowToObject func(*sql.Rows, *T) error, 
	queryString string,
	args ...string,
) ([]T, error) {

	rows, queryError := client.Query(queryString)

	var empty []T

	if queryError != nil {
		return empty, queryError 	
	}

	capacity := 100

	rowArray := make([]T, 0, capacity)
	var index int

	for rows.Next() {
		var receiverObject *T

		if index == capacity - 1 {
			capacity += 100
			newRowArray := make([]T, 0, capacity)
			
			copy(newRowArray, rowArray)
			rowArray = newRowArray
		}

		// scanError := rows.Scan()

		scanError := scanRowToObject(rows, receiverObject)

		if scanError != nil {
			return empty, scanError
		}

		rowArray[index] = *receiverObject
		index++
	}

	return rowArray, nil
}
