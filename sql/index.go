package sql

import (
	"database/sql"
	"log"
	"strconv"
	"errors"
	"strings"
	"fmt"
	"reflect"

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

func ReceiveRows[T SQLMetaStruct](rows *sql.Rows) ([]T, error) {

	var empty []T

	capacity := 100

	rowArray := make([]T, capacity, capacity)
	var index int

	for rows.Next() {

		var receiver T 
		zeroedStruct := receiver.Zero()

		asserted, ok := zeroedStruct.(T)

		if !ok {
			log.Printf("type: %T", zeroedStruct)
			return nil, errors.New("type assertion failed")
		}

		receiver = asserted

		if index == capacity-1 {
			capacity += 100
			newRowArray := make([]T, 0, capacity)

			copy(newRowArray, rowArray)
			rowArray = newRowArray
		}


		err := ScanRow[T](rows, receiver)

		if err != nil {
			log.Println("scanError", err.Error())
			return empty, err
		}

		rowArray[index] = receiver
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
func QueryForStructs[T SQLMetaStruct](
	client PGClient,
	scanRowToObject func(*sql.Rows, T) error,
	queryString string,
	args ...any,
) ([]T, error) {

	var empty []T

	rows, queryError := client.Query(queryString, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[T](rows)
}

// func SimpleQuery(
// 	client PGClient,
// 	queryString string,
// 	args ...any,
// ) error {

// 	rows, queryError := client.Query(queryString, args...)

// 	if queryError != nil {
// 		return queryError
// 	}

// 	_, receiveError := ReceiveRows[SQLReporter](rows, ScanRowNoOp[any])

// 	return receiveError
// }

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
		args[i] = arg
	}
	return args
}

type SQLReporter interface {
	Keys() []string
	Values() []any
	Get(string) (any, error)
	Set(string, any) error
	TableName() string
	Init() SQLReporter
}

// if a `nil` is passed in `[]S` this crashes.
func InsertStructs[S SQLMetaStruct](client PGClient, rows []S) error {

	rowCount := len(rows)

	if rowCount == 0 {
		return errors.New("no rows, nothing to insert")
	}

	// handle a panic, possible if a value in rows is nil
    defer func() {
        if panicVal := recover(); panicVal != nil {
			log.Println("InsertStructs() panic value:", panicVal)
        }
    }()

	tableName := rows[0].TableName()

	if tableName == "" {
		return errors.New("all rows are nil, nothing to insert")
	}

	var empty S

	keys := empty.Keys()

	columnNamesString := strings.Join(empty.Keys(), ", ")
	columnCount := len(keys)

	query := strings.Builder{}

	query.WriteString("INSERT INTO " + tableName)
	query.WriteString("(" + columnNamesString + ")")
	query.WriteString(" VALUES ")

	size := rowCount * columnCount
	// log.Println("rowCount:", rowCount)
	// log.Println("columnCount:", columnCount)
	// log.Println("size:", size)
	values := make([]any, 0, size)

	seq := SQLArgSequence{}

	// var topIndex int
	for index, row := range rows {
		
		last := index == rowCount - 1

		values = append(values, row.Values()...)

		nextArgs := seq.NextNString(columnCount)

		// log.Println("nextArgs:", nextArgs)

		query.WriteString("(")
		query.WriteString(strings.Join(nextArgs, ", "))
		query.WriteString(")")

		if !last {
			query.WriteString(",")
		}
	}

	query.WriteString(";")

	queryString := query.String()

	log.Println("queryString", queryString)

	_, err := client.Query(queryString, values...)

	if err != nil {
		log.Println("error running query:", queryString)
		return err
	}

	return nil
}

func ScanRow[T SQLMetaStruct](rows *sql.Rows, object T) error {

	// if object == nil {
	// 	return errors.New("ScanRow received a nil pointer")
	// }

	fmt.Printf("OBJECT: %+v\n", object)
	// log.Println("object", object)

	values := object.Values()

	fmt.Printf("VALUES: %+v\n", values)

	for _, value := range values {
		fmt.Printf("VALUE: %+v\n", value)
		fmt.Println("VALUE type:", reflect.TypeOf(value))
	}


	// log.Println("values", values)

	// for _, v := range values {
	// 	fmt.Println("Indirect type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Type()) // prints main.CustomStruct

	// 	fmt.Println("Indirect value type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Kind()) // prints struct
	// }

	err := rows.Scan(values...)

	if err != nil {
		log.Println("SCAN ERRUR, BRUH")
		return err
	}

	return nil
}

// if a `nil` is passed in `[]S` this crashes.
func GetStructs[S SQLMetaStruct](client PGClient, query string, args []any) ([]S, error) {

	var empty []S

	rows, queryError := client.Query(query, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[S](rows)
}

func MetaQuery[S SQLMetaStruct](client PGClient, query string, args []any) ([]S, error) {

	var empty []S

	rows, queryError := client.Query(query, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[S](rows)
}
