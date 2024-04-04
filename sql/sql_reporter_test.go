package sql

import (
	// "database/sql"
	// "errors"
	"log"
	//"strconv"
	"testing"
)

func TestInsertStructsWithSQLMetaStruct(t *testing.T) {

	db, err := BuildPostgresClient("user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil || db == nil {
		t.Fatal("error getting client during TestInsertStructs()")
	}

	zero := 0
	one := 1
	two := 2
	three := float32(0)
	four := "blah blah blah test"
	five := "2024-03-25 01:58:08.789206+00"
	six := "2024-03-25 01:58:08.789206+00"
	seven := "2024-03-25 01:58:08.789206+00"

	users := []Expenditure{ Expenditure{ &zero, &one, &two, &three, &four, &five, &six, &seven } }

	log.Println("inserting...")

	err = InsertStructs[Expenditure](db, users)

	if err != nil {
		t.Fatal(err)
	}

	log.Println("now selecting...")

	query := "select user_id, category_id, value, description, date_occurred, create_date, modified_date from expenditure;"

	newRows, err := MetaQuery[Expenditure](db, query, nil)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("NEW ROWS:", newRows)

	for i, newRow := range newRows {

		newName := *newRow.Description
		oldName := *users[i].Description

		if newName != oldName {
			t.Fatal(err)
		}
	}
}

func TestInsertVanillaAny(t *testing.T) {

	t.Skip()

	client, err := BuildPostgresClient("user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil || client == nil {
		t.Fatal("error getting client during TestInsertVanillaAny()")
	}

	zero := 0
	one := 1
	two := 2
	three := float32(0)
	four := "blah blah blah test"
	five := "2024-03-25 01:58:08.789206+00"
	six := "2024-03-25 01:58:08.789206+00"
	seven := "2024-03-25 01:58:08.789206+00"

	users := []Expenditure{ Expenditure{ &zero, &one, &two, &three, &four, &five, &six, &seven } }

	log.Println("inserting...")

	err = InsertStructs[Expenditure](client, users)

	if err != nil {
		t.Fatal(err)
	}

	query := "select user_id, category_id, value, description, date_occurred, create_date, modified_date from expenditure;"

	args := make([]any, 0)

	rows, err := client.Query(query, args...)

	if err != nil {
		t.Fatal(err)
	}

	results := [][]any{}

	for rows.Next() {

		two := 0
		three := 0
		four := float32(0)
		five := ""
		six := ""
		seven := ""
		eight := ""

		values := []any{ &two, &three, &four, &five, &six, &seven, &eight }

		err := rows.Scan(values...)

		if err != nil {
			t.Fatal(err)
		}
		results = append(results, values)
	}

	t.Log("before, results length:", len(results))

	for i, values := range results {
		t.Log("i:", i)
		for n, value := range values {
			t.Log("i, n", i, n, "value", value)
		}
	}

	t.Log("done.")
}