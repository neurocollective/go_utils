package go_utils

import (
	"testing"
	"database/sql"
	"log"
	"errors"
)

type TestStruct struct {
	id int
	name string
}

func ScanForTestStruct(rows *sql.Rows, tester *TestStruct) error {

	if rows == nil {
		return errors.New("rows is nil inside ScanForTestStruct")
	}

	// if tester == nil {
	// 	tester = new(TestStruct)
	// }

	id := tester.id
	name := tester.name

	scanError := rows.Scan(&id, &name)
	
	if scanError != nil {
		return scanError
	}

	log.Println("tester id:", tester.id)
	log.Println("tester name:", tester.name)

	return nil
}

func TestQueryForStructs(t *testing.T) {

	// db, getClientError := BuildPostgresClient("postgresql://postgres:postgres@localhost:5432/postgres")

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

	t.Log("testStructs", testStructs)
}