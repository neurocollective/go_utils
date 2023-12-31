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
	if err != nil {

        log.Println("ERROR opening postgres connection with github.com/neurocollective/go_utils.BuildPostgresClient() ->")
        log.Println(err.Error())

		return nil, err
	}

	return db, nil
}

func ReceiveRows[T any](rows *sql.Rows, scanRowToObject func(*sql.Rows, *T) error) ([]T, error) {

	var empty []T

	capacity := 100

	rowArray := make([]T, capacity, capacity)
	var index int

	for rows.Next() {

		receiverObject := new(T)

		if index == capacity - 1 {
			capacity += 100
			newRowArray := make([]T, 0, capacity)
			
			copy(newRowArray, rowArray)
			rowArray = newRowArray
		}

		scanError := scanRowToObject(rows, receiverObject)

		if scanError != nil {
            log.Println("scanError", scanError.Error())
			return empty, scanError
		}

		rowArray[index] = *receiverObject
		index++
	}

    getNextRowError := rows.Err()

    if getNextRowError != nil {
        log.Println("error getting next row:", getNextRowError.Error())
        return empty, getNextRowError
    }

	return rowArray[:index], nil

}

// takes a struct-specific `scanRows`
func QueryForStructs[T any](
	client PGClient, 
	scanRowToObject func(*sql.Rows, *T) error, 
	queryString string,
	args ...any,
) ([]T, error) {

	var empty []T

	rows, queryError := client.Query(queryString, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[T](rows, scanRowToObject)
}
