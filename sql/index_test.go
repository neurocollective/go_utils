 package go_utils

import (
	"testing"
	"database/sql"
	"log"
	"errors"
)

type TestStruct struct {
	id int `nc:"id"`
	name string `nc:"name"`
}

func (t TestStruct) GetStructKeys() []string {
	return []string { "id", "name" }
}

func GetFromTestStruct[T any](testStruct TestStruct, structFieldName string) (T, error) {
	switch structFieldName {
		case "id":
			idValue, ok := (testStruct.id).T
			if !ok {
				return idValue, errors.New("id cannot be retrieved as requested type.")
			}
			return idValue, nil
		case "name":
			nameValue, ok := (testStruct.name).T
			if !ok {
				return nameValue, errors.New("name cannot be retrieved as requested type.")
			}
			return nameValue, nil
		default:
			return nil, errors.New(structFieldName + " not a valid struct field")
	}
}

func (t TestStruct) Get(structFieldName string) (any, error) {
	switch structFieldName {
		case "id":
			return t.id, nil
		case "name":
			return t.name, nil
		default:
			return nil, errors.New(structFieldName + " not a valid struct field")
	}
}

func ScanForTestStruct(rows *sql.Rows, tester *TestStruct) error {

	log.Println("tester inside ScanForTestStruct:", tester)

	if rows == nil {
		return errors.New("rows is nil inside ScanForTestStruct")
	}

	// if tester == nil {
	// 	tester = new(TestStruct)
	// }

	idPointer := &tester.id
	namePointer := &tester.name

	scanError := rows.Scan(idPointer, namePointer)

	log.Println("idPointer inside ScanForTestStruct:", *idPointer)
	log.Println("namePointer inside ScanForTestStruct:", *namePointer)

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