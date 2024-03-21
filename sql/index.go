package go_utils

import (
	"database/sql"
	"log"
	"strconv"
	"errors"

	_ "github.com/lib/pq"
)

type PGClient interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// connectionString -> "user=postgres password=postgres dbname=postgres sslmode=disable"
func BuildPostgresClient(connectionString string) (PGClient, error) {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {

		log.Println("ERROR opening postgres connection with github.com/neurocollective/go_utils.BuildPostgresClient() ->")
		log.Println(err.Error())

		return nil, err
	}

	return db, nil
}

func ScanRowNoOp[T any](rows *sql.Rows, object *T) error {
	log.Println(rows)
	log.Println(object)
	return nil
}

func ReceiveRows[T any](rows *sql.Rows, scanRowToObject func(*sql.Rows, *T) error) ([]T, error) {

	var empty []T

	capacity := 100

	rowArray := make([]T, capacity, capacity)
	var index int

	for rows.Next() {

		receiverObject := new(T)

		if index == capacity-1 {
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

func SimpleQuery(
	client PGClient,
	queryString string,
	args ...any,
) error {

	rows, queryError := client.Query(queryString, args...)

	if queryError != nil {
		return queryError
	}

	_, receiveError := ReceiveRows[any](rows, ScanRowNoOp[any])

	return receiveError
}

type SQLArgSequence struct {
	Id int
}

func (as *SQLArgSequence) Next() int {
	currentId := as.Id

	if currentId == 0 {
		as.Id += 1
		currentId = 1
	}

	as.Id += 1
	return currentId
}

func (as *SQLArgSequence) NextN(n int) []int {

	ids := make([]int, n, n)

	for i := 0; i < n; i++ {
		id := as.Next()
		ids = append(ids, id)
	}
	return ids
}

func (as *SQLArgSequence) NextString() string {
	id := as.Next()

	return "$" + strconv.Itoa(id)
}

func (as *SQLArgSequence) NextNString(n int) []string {
	args := make([]string, n, n)

	for i := 0; i < n; i++ {
		arg := as.NextString()
		args = append(args, arg)
	}
	return args
}

type SQLReporter interface {
	Keys() []string
	Values() []any
	Get(string) (error, any)
	Set(string, any) error
	TableName() string
}

// if a `nil` is passed in `[]S` this crashes.
func InsertStructs[S SQLReporter](client PGClient, rows []S) error {

	var tableName string

	rowCount := len(rows)

	if rowCount == 0 {
		return errors.New("no rows, nothing to insert")
	}

	// handle a panic, possible if a value in rows is nil
    defer func() {
        if panicVal := recover(); panicVal != nil {
			log.Println("InsertStructs() panic value:", panicVal)
			return errors.New("nil passed to InsertStructs()")
        }
    }()

	// will panic if nil
	for _, row := range rows {
		tableName := row.TableName(0)
	}

	if tableName == "" {
		return errors.New("all rows are nil, nothing to insert")
	}

	var empty S

	columnNames := strings.Join(empty.Keys(), ", ")
	columnCount := len(columnNames)

	query := strings.Builder{}

	query.WriteString("INSERT INTO " + tableName)
	query.WriteString("(" + columnNames + ")")
	query.WriteString(+ "VALUES ")

	size := rowCount * columnCount

	values := make([]any, size, size)

	seq := SQLArgSequence{}

	for index, row := range rows {

		last := index == rowCount - 1

		column := columnNames[index]
		values = append(values, row.Values()...)

		nextArgs := seq.NextNString(columnCount)

		query.WriteString("(")
		query.WriteString(strings.Join(nextArgs, ", "))
		query.WriteString(")")

		if !last {
			query.WriteString(",")
		}
	}

	query.WriteString(";")

}
