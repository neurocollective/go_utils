package generator

import (
	"os"
	"log"
	"strings"
)

type GenerationConfig struct {
	Fields []map[string]string // { "fieldName": "id", "type": "int" }
	StructName string
	OutputPath string
	TemplatePath string
	PackageName string
}

const switchCaseEntryLineCount = 6

func Generate(config GenerationConfig) {

	fileBytes, readError := os.ReadFile(config.TemplatePath)

	if readError != nil {
		os.Exit(1)
	}

	fileString := string(fileBytes)

	fileLines := strings.Split(fileString, "\n")

	formattedLines := make([]string, len(fileLines))

	// var funcFound bool

	for index, line := range fileLines {
		log.Println(line)

		if index == 0 {
			hasPackage := strings.Contains(line, "package")

			if !hasPackage {
				log.Println("first line of template does not have a package declaration")
				os.Exit(1)
			}
			formattedLines[index] = strings.Replace(line, "$packageName", config.PackageName, -1)
			continue
		}

		// if !funcFound {
		// 	funcFound := 
		// }
	}

	outputBytes := []byte(strings.Join(formattedLines, "\n"))

	writeError := os.WriteFile(config.OutputPath, outputBytes, 0666)

	if writeError != nil {
		os.Exit(1)
	}
}