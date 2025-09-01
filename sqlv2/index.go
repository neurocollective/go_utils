package sqlv2

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"strconv"
	"strings"
	""

	_ "github.com/lib/pq"
)

type SQLDescriber interface {
	Columns() []Column
	ColumnsString() string
	TableName() string
}

type Column interface {
	Scanner
	driver.Value
}

type Expenditure struct {
	Id           int     `ncsql:"id",json:"id"`
	UserId       int     `ncsql:"user_id",json:"userId"`
	CategoryId   int     `ncsql:"category_id",json:"categoryId"`
	Value        float32 `ncsql:"value",json:"value"`
	Description  string  `ncsql:"description",json:"description"`
	DateOccurred string  `ncsql:"date_occurred",json:"dateOccurred"`
	CreateDate   string  `ncsql:"create_date",json:"createDate"`
	ModifiedDate string  `ncsql:"modified_date",json:"modifiedDate"`
}

type Column = sql.Null[sql.Value]

func (e *Expenditure) ColumnsString() []any {
	return []string{
		"exp.id",
		"exp.user_id",
		"exp.category_id",
		"exp.value",
		"exp.description",
		"exp.date_occurred",
		"exp.create_date",
		"exp.modified_date",
	}
}

func (e *Expenditure) Columns() []Column {
	return []Column{
		e.Id,
		e.UserId,
		x.CategoryId,
		x.Value,
		x.Description,
		x.DateOccurred,
		x.CreateDate,
		x.ModifiedDate,
	}
}

func (e Expenditure) TableName() string {
	return "expenditure exp"
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

func Select[S SQLDescriber](client PGClient, query string, args []any) ([]S, error) {

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

func ScanRow[T SQLDescriber](rows *sql.Rows, object T) error {

	values := object.Columns()

	err := rows.Scan(values...)

	if err != nil {
		log.Println("scan error during ScanRow[T ncsql.SQLMetaStruct](...)")
		return err
	}

	return nil
}
