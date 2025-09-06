package main

import (
	"database/sql"
	"database/sql/driver"
	//"errors"
	"log"
	//"strconv"
	//"strings"

	_ "github.com/lib/pq"
)

type SQLDescriber interface {
	Columns() []any
	ColumnsString() string
	TableName() string
}

// type Column = sql.Scanner

// `NotNull[T]` implements `sql.Scanner`
// `driver.Value` is just `any`, but should be certain pointer types.
// https://pkg.go.dev/database/sql/driver#Value
type NotNull[T driver.Value] struct {
	V     T
}

func (nn *NotNull[T]) Scan(value any) error {

	// if value == nil {
	// 	return errors.New("nil value for NotNull Column")
	// }

	// nullable := &sql.Null[T]{nn.V,true}
	nullable := new(sql.Null[T])
	err := nullable.Scan(value)

	if err != nil {
		log.Println("err:", err)
	} else {
		nn.V = nullable.V
	}

	return err
}

func (nn NotNull[T]) Value() (driver.Value, error) {
	nullable := sql.Null[T]{nn.V,true}
	return nullable.Value()
}

// func (n Null[T]) Value() (driver.Value, error)

// connectionString -> "user=postgres password=postgres dbname=postgres sslmode=disable"
// func BuildPostgresClient(connectionString string) (PGClient, error) {

// 	db, err := sql.Open("postgres", connectionString)
// 	if err != nil {

// 		log.Println("ERROR opening postgres connection with github.com/neurocollective/go_utils.BuildPostgresClient() ->")
// 		log.Println(err.Error())

// 		return nil, err
// 	}

// 	return db, nil
// }

func Select[S SQLDescriber](client *sql.DB, query string, args []any) ([]S, error) {

	var empty []S

	rows, queryError := client.Query(query, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[S](rows)
}

func ReceiveRows[T SQLDescriber](rows *sql.Rows) ([]T, error) {

	var empty []T

	capacity := 100

	rowArray := make([]T, capacity, capacity)
	var index int

	for rows.Next() {

		var receiver T

		if index == capacity-1 {
			capacity += 100
			newRowArray := make([]T, 0, capacity)

			copy(newRowArray, rowArray)
			rowArray = newRowArray
		}

		values := receiver.Columns()

		err := rows.Scan(values...)

		if err != nil {
			log.Println("scanError", err.Error())
			return empty, err
		}

		rowArray[index] = receiver
		index++
	}

	// getNextRowError := rows.Err()

	// if getNextRowError != nil {
	// 	log.Println("error getting next row:", getNextRowError.Error())
	// 	return empty, getNextRowError
	// }

	return rowArray[:index], nil
}
