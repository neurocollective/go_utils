package sql

import (
	"log"
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

	expenditures := []Expenditure{ Expenditure{ &zero, &one, &two, &three, &four, &five, &six, &seven } }

	log.Println("inserting...")

	err = Insert[Expenditure](db, expenditures)

	if err != nil {
		t.Fatal(err)
	}

	log.Println("now selecting...")

	selectQuery := "select id, user_id, category_id, value, description, date_occurred, create_date, modified_date from expenditure;"

	newRows, err := Select[Expenditure](db, selectQuery, nil)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("NEW ROWS:", newRows)

	for i, newRow := range newRows {

		newName := *newRow.Description
		oldName := *expenditures[i].Description

		if newName != oldName {
			t.Fatal(err)
		}
	}

	row := newRows[0]

	testDescription := "test"

	row.Description = &testDescription

	err = Update[Expenditure](db, row)

	if err != nil {
		t.Fatal(err)
	}

	newRows, err = Select[Expenditure](db, selectQuery, nil)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("UPDATED ROWS:", newRows)

	if *newRows[0].Description != testDescription {
		t.Fatalf("unexpected row.Description, expected '%s' but got '%s'", testDescription, *newRows[0].Description)
	}

}
