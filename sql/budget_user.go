package go_utils

import (
	"log"
	"errors"
)

// a valid SQLReporter will need pointers as every field. `SQLReporter` implies this.
type BudgetUser struct {
	Name *string `ncsql:"name",json:"name"`
}

func (bu *BudgetUser) Init() SQLReporter {

	name := ""

	new := new(BudgetUser)
	new.Name = &name
	return new
}

func (bu *BudgetUser) Keys() []string {
	return []string{
		"name",
	}
}

func (bu *BudgetUser) TableName() string {
	return "budget_user"
}

func (bu *BudgetUser) Values() []any {

	log.Println("bu in Values()", bu)

	return []any{
		bu.Name,
	}
}

func (bu *BudgetUser) Get(structKey string) (any, error) {
	values := bu.Values()
	for index, key := range bu.Keys() {
		value := values[index]
		if key == structKey {
			return value, nil
		}
	}
	return nil, errors.New(structKey + "not found")
}

func (bu *BudgetUser) Set(structKey string, newValue any) error {
	values := bu.Values()
	for index, key := range bu.Keys() {
		if key == structKey {
			values[index] = newValue
		}
		return nil
	}
	return errors.New(structKey + "not found")
}