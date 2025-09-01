package main

import (
	"fmt"
	"database/sql"
)

func BuildPostgresClient(connectionString string) (sql.DB, error) {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {

		log.Println("ERROR opening postgres connection with github.com/neurocollective/go_utils.BuildPostgresClient() ->")
		log.Println(err.Error())

		return nil, err
	}

	return db, nil
}

func main() {

	client, err := BuildPostgresClient("user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		fmt.Println("error connecting to be", err)
		os.Exit(1)
	}

	rows, queryError := client.Query(queryString, args...)

	if queryError != nil {
		return empty, queryError
	}

}
