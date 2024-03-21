package go_utils

import (
	"log"
	"errors"
)

// a valid SQLReporter will need pointers as every field. `SQLReporter` implies this.
type Person struct {
	Name *string `ncsql:"name",json:"name"`
	Greeting *string `ncsql:"greeting",json:"greeting"`
}

func (p *Person) SayHello() string {
	return "My name is " + (*(p.Name)) + ", " + (*(p.Greeting))
}

func (p *Person) Keys() []any {
	return []any{
		"name",
		"greeting",
	}
}

func (p *Person) Values() []any {
	return []any{
		p.Name,
		p.Greeting,
	}
}

func (p *Person) Get(structKey string) (any, error) {
	values := p.Values()
	for index, key := range p.Keys() {
		value := values[index]
		if key == structKey {
			return value, nil
		}
	}
	return nil, errors.New(structKey + "not found")
}

func (p *Person) Set(structKey string, newValue any) error {
	values := p.Values()
	for index, key := range p.Keys() {
		if key == structKey {
			values[index] = newValue
		}
		return nil
	}
	return errors.New(structKey + "not found")
}