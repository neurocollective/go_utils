package go_utils

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

	query := "select name from expenditure;"

	newRows, err := GetStructs[Expenditure](db, query, nil)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(newRows)

	// for i, newRow := range newRows {

	// 	newName := *newRow.Name
	// 	oldName := *users[i].Name

	// 	if newName != oldName {
	// 		t.Fatal(err)
	// 	}
	// }
}