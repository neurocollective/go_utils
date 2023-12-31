package main

import (
	"github.com/neurocollective/go_utils/generator"
	"log"
	"os"
)

func main() {

	cwd, getCwdError := os.Getwd()

	if getCwdError != nil {
		os.Exit(1)
	}

	log.Println("cwd:", cwd)

	fields := []map[string]string{
		map[string]string{"fieldName": "id", "type": "int"},
		map[string]string{"fieldName": "name", "type": "string"},
	}

	config := generator.GenerationConfig{
		fields,
		"TestStruct",
		cwd + "/generated",
		cwd + "/generator/index.templ",
		"generatorTest",
	}

	generator.Generate(config)
}
