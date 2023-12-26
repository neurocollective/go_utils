package go_utils

import (
	"testing"
	"database/sql"
	"log"
)

type TestStruct struct {
	id int
	name string
}

type FakeDb struct {
	pointlessField string
}

func (f FakeDb)Query(query string, args ...any) (*sql.Rows, error){
	log.Println(query)
	log.Println(args...)
	return nil, nil
}

func TestQueryForStructs(t *testing.T) {

	id := 1

	scan := func(rows *sql.Rows, tester *TestStruct) error {
		log.Println(rows)
		log.Println(tester)

		tester.id = id;
		tester.name = "Some dude"
		id++
		return nil
	}

	query := "SELECT * from test_table where id = $1;"

	args := []string{ "1" }

	db := FakeDb{}

	_, parseError := QueryForStructs[TestStruct](db, scan, query, args...)

	if parseError != nil {
		t.Fatal("error!" + parseError.Error())
	}
}