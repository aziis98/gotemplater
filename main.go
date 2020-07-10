package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
)

var helpMessage = strings.TrimLeftFunc(`
gotemplater - A super small CLI utility wrapping template.Execute()
              and template.ExecuteTemplate()

USAGE: 
gotemplater [OPTIONS...] [FILES...] < [CONFIG_FILE...]

OPTIONS:
    -h, --help		Shows this message
    
    -o, --output FILE	File to write to, by default uses stdout
    
    -e, --execute NAME 	Template name to pass to .ExecuteTemplate(), otherwise uses .Execute()

    -d, --data DATA_FILE    Data file to use, otherwise use stdin
    -f, --format FORMAT	    Format for the data file, can be one of: *"json", "yaml"
                            *the default format is JSON.

EXAMPLES:

    Print to stdout: 
        gotemplater template.html < data.json

    Print to file: 
        gotemplater -o rendered.html template.html < data.json

    Get data from file in another format: 
        gotemplater -f yaml -d data.yaml -o rendered.html template.html

`, unicode.IsSpace)

func main() {
	args := os.Args[1:]
	templateFiles := []string{}

	execute := ""
	outputFile := ""

	dataFormat := "json"
	dataFile := ""

	if len(args) == 0 {
		fmt.Print(helpMessage)
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
		case "-h", "--help":
			fmt.Print(helpMessage)
			os.Exit(0)
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
			dataFormat = strings.ToLower(flagValue)
			args = args[2:]
		default:
			if strings.HasPrefix(flag, "-") {
				log.Fatal(fmt.Errorf(`Unrecognized flag "%s"`, flag))
			}

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
