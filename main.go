package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const helpMessage = `
gotemplater - CLI utility to render Go Templates to files
`

func main() {
	args := os.Args[1:]
	templateFiles := []string{}

	execute := ""
	outputFile := ""

	dataFormat := "json"
	dataFile := ""

	if len(args) == 0 {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	for {
		if len(args) == 0 {
			break
		}

		var flag, flagValue, fileName string
		if len(args) > 1 {
			flagValue = args[1]
		}
		{
			flag = args[0]
			fileName = args[0]
		}

		switch flag {
		case "-e", "--execute":
			execute = flagValue
			args = args[2:]
		case "-o", "--output":
			outputFile = flagValue
			args = args[2:]
		case "-d", "--data":
			dataFile = flagValue
			args = args[2:]
		case "-f", "--format":
			dataFormat = flagValue
			args = args[2:]
		default:
			templateFiles = append(templateFiles, fileName)
			args = args[1:]
		}
	}

	var bytes []byte

	if dataFile == "" {
		var err error
		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		bytes, err = ioutil.ReadFile(dataFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var data map[string]interface{}
	switch dataFormat {
	case "json":
		err := json.Unmarshal(bytes, &data)
		if err != nil {
			log.Fatal(err)
		}
	case "yaml":
		err := yaml.Unmarshal(bytes, &data)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(fmt.Errorf(`The data format "%s" is invalid or not supported, use "json" or "yaml"`, dataFormat))
	}

	var output io.Writer
	if outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}

		output = file
	}

	tmpl, err := template.New("").ParseFiles(templateFiles...)
	if err != nil {
		log.Fatal(err)
	}

	if execute == "" {
		tmpl.Execute(output, data)
	} else {
		tmpl.ExecuteTemplate(output, execute, data)
	}
}
