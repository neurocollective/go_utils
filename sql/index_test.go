 package go_utils

import (
	"testing"
	"database/sql"
	"log"
	"errors"
	"strconv"
)

type TestStruct struct {
	id int `nc:"id"`
	name string `nc:"name"`
}

func (t TestStruct) GetStructKeys() []string {
	return []string { "id", "name" }
}

func ScanForTestStruct(rows *sql.Rows, tester *TestStruct) error {

	if rows == nil {
		return errors.New("rows is nil inside ScanForTestStruct")
	}

	idPointer := &tester.id
	namePointer := &tester.name

	scanError := rows.Scan(idPointer, namePointer)

	if scanError != nil {
		return scanError
	}

	return nil
}

func ScanForTestStructIds(rows *sql.Rows, tester *TestStruct) error {
	if rows == nil {
		return errors.New("rows is nil inside ScanForTestStruct")
	}

	idPointer := &tester.id

	scanError := rows.Scan(idPointer)

	if scanError != nil {
		return scanError
	}

	return nil
}

func TestQueryForStructs(t *testing.T) {

	db, getClientError := BuildPostgresClient("user=postgres password=postgres dbname=postgres sslmode=disable")

	if getClientError != nil || db == nil {
		t.Fatal("error getting client")
	}

	// query := "SELECT id, name from budget_user where id = $1;"
	// args := []any{ 1 }

	query := "SELECT id, name from budget_user where id = 1;"
	args := []any{}

	testStructs, parseError := QueryForStructs[TestStruct](db, ScanForTestStruct, query, args...)

	if parseError != nil {
		t.Fatal("error!", parseError.Error())
	}

	if len(testStructs) == 0 {
		t.Fatal("no results in array!")		
	}

	if len(testStructs) > 1 {
		t.Fatal("too many results in array!")		
	}

	receivedId := testStructs[0].id
	receivedName := testStructs[0].name

	if receivedId != 1 {
		t.Fatal("did not receive expected id of 1! instead got", receivedId)
	}

	if receivedName != "david" {
		t.Fatal("did not receive expected name of \"dave\"! instead got", receivedName)
	}
}

type TestExpenditure struct {
	Id int
	UserId int
	CategoryId *int
	Value float32
	Description string
	DateOccurred string
}

func TestInsertStructsQuery(t *testing.T) {

	expenditures := []TestExpenditure{
		TestExpenditure{ -1, 1, nil, 20.99, "Blockchain Backscratcher", "2023-12-30 21:49:01.172639+00" },
		TestExpenditure{ -1, 1, nil, 900.00, "NOT Cocaine", "2023-12-30 21:49:01.172639+00" },
		TestExpenditure{ -1, 1, nil, 4000.00, "Darkweb hitman", "2023-12-30 21:49:01.172639+00" },
	}

	db, getClientError := BuildPostgresClient("user=postgres password=postgres dbname=postgres sslmode=disable")		

	if getClientError != nil || db == nil {
		t.Fatal("error getting client")
	}

	queryStem := "insert into expenditure (id, user_id, category_id, value, description, date_occurred) values "

	var queryValuesSuffix string
	size := len(expenditures) * 5
	args := make([]any, size, size)

	argumentIndex := 1

	for index, ex := range expenditures {

		first := " $" + strconv.Itoa(argumentIndex) + ", "
		args[argumentIndex - 1] = ex.UserId
		argumentIndex++

		second := "$" + strconv.Itoa(argumentIndex) + ", "
		args[argumentIndex - 1] = ex.CategoryId
		argumentIndex++

		third := "$" + strconv.Itoa(argumentIndex) + ", "
		args[argumentIndex - 1] = ex.Value
		argumentIndex++

		fourth := "$" + strconv.Itoa(argumentIndex) + ", "
		args[argumentIndex - 1] = ex.Description
		argumentIndex++

		fifth := "$" + strconv.Itoa(argumentIndex)
		args[argumentIndex - 1] = ex.DateOccurred
		argumentIndex++

		queryValues := "(nextval('expenditure_id_seq'),"

		if index == len(expenditures) -1 {
			queryValues += first + second + third + fourth + fifth + ")"			
		} else {
			queryValues += first + second + third + fourth + fifth + "), "
		}
		queryValuesSuffix += queryValues
	}

	fullQuery := queryStem + queryValuesSuffix + " RETURNING id;"

	log.Println("fullQuery:", fullQuery)
	log.Println("args:", args)

	testExpenditures, parseError := QueryForStructs[TestStruct](db, ScanForTestStructIds, fullQuery, args...)

	if parseError != nil {
		t.Fatal("error!", parseError.Error())
	}

	log.Println(testExpenditures)


}