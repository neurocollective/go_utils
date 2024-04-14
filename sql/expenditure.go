package sql

import (
	"errors"
)

// a valid SQLReporter will need pointers as every field. `SQLReporter` implies this.
type Expenditure struct {
	Id           *int     `ncsql:"id",json:"id"`
	UserId       *int     `ncsql:"user_id",json:"userId"`
	CategoryId   *int     `ncsql:"category_id",json:"categoryId"`
	Value        *float32 `ncsql:"value",json:"value"`
	Description  *string  `ncsql:"description",json:"description"`
	DateOccurred *string  `ncsql:"date_occurred",json:"dateOccurred"`
	CreateDate   *string  `ncsql:"create_date",json:"createDate"`
	ModifiedDate *string  `ncsql:"modified_date",json:"modifiedDate"`
}

func (e Expenditure) GetId() *int {
	return e.Id
}

func (e Expenditure) Zero() SQLMetaStruct {

	new := Expenditure{}

	one := 0
	two := 0
	three := 0
	four := float32(0)
	five := ""
	six := ""
	seven := ""
	eight := ""

	new.Id = &one
	new.UserId = &two
	new.CategoryId = &three
	new.Value = &four
	new.Description = &five
	new.DateOccurred = &six
	new.CreateDate = &seven
	new.ModifiedDate = &eight

	return new
}

// this should be generated code, based on column names
func (e Expenditure) Keys() []string {
	return []string{
		"user_id",
		"category_id",
		"value",
		"description",
		"date_occurred",
		"create_date",
		"modified_date",
	}
}

func (e Expenditure) KeysAll() []string {
	return []string{
		"id",
		"user_id",
		"category_id",
		"value",
		"description",
		"date_occurred",
		"create_date",
		"modified_date",
	}
}

func (e Expenditure) TableName() string {
	return "expenditure"
}

func (e Expenditure) Values() []any {

	return []any{
		e.UserId,
		e.CategoryId,
		e.Value,
		e.Description,
		e.DateOccurred,
		e.CreateDate,
		e.ModifiedDate,
	}
}

func (e Expenditure) ValuesAll() []any {

	return []any{
		e.Id,
		e.UserId,
		e.CategoryId,
		e.Value,
		e.Description,
		e.DateOccurred,
		e.CreateDate,
		e.ModifiedDate,
	}
}

func (e Expenditure) Get(structKey string) (any, error) {
	values := e.Values()
	for index, key := range e.Keys() {
		value := values[index]
		if key == structKey {
			return value, nil
		}
	}
	return nil, errors.New(structKey + "not found")
}
