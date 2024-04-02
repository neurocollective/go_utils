package sql

import (
	"database/sql"
	"log"
	"strconv"
	"errors"
	"strings"

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

func ReceiveRows[T SQLMetaStruct](rows *sql.Rows, scanRowToObject func(*sql.Rows, T) error) ([]T, error) {

	var empty []T

	capacity := 100

	rowArray := make([]T, capacity, capacity)
	var index int

	for rows.Next() {

		// receiverObject := new(T)
		var receiverObject T
		// receiverPointer := &receiverObject

		log.Println("receiverObject:", receiverObject)
		// log.Println("receiverPointer:", receiverPointer)

		receiverPointer := receiverObject.Init()

		log.Println("receiverPointer:", receiverPointer)

		for _, value := range receiverPointer.Values() {
			log.Println(value)
		}

		if index == capacity-1 {
			capacity += 100
			newRowArray := make([]T, 0, capacity)

			copy(newRowArray, rowArray)
			rowArray = newRowArray
		}

		assertedType, ok := receiverPointer.(T)

		if !ok {
			 log.Printf("got data of type %T, not the asserted value:", receiverPointer)
			return []T{}, errors.New("I don't even know")
		}

		scanError := scanRowToObject(rows, assertedType)

		if scanError != nil {
			log.Println("scanError", scanError.Error())
			return empty, scanError
		}

		rowArray[index] = assertedType
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

	return ReceiveRows[T](rows, scanRowToObject)
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

	log.Println("object", object)

	values := object.Values()

	log.Println("values", values)

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

	return ReceiveRows[S](rows, ScanRow)
}

func MetaQuery[S SQLMetaStruct](client PGClient, query string, args []any) ([]S, error) {

	var empty []S

	rows, queryError := client.Query(query, args...)

	if queryError != nil {
		return empty, queryError
	}

	return ReceiveRows[S](rows, ScanRow)
}
