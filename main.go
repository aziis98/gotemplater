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
    -h, --help			Shows this message
    
    -o, --output FILE		File to write to, by default uses stdout
    
    -e, --execute NAME		Template name to pass to .ExecuteTemplate(), otherwise uses .Execute()

    -c, --content FILE		Adds a "content"/"Content" variable to the context of the template for 
				rendering importing data from files

    -d, --data DATA_FILE	Data file to use, otherwise use stdin
    -f, --format FORMAT		Format for the data file, can be one of: *"json", "yaml"
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
	contentFile := ""
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
		flag = args[0]
		fileName = args[0]

		switch flag {
		case "-h", "--help":
			fmt.Print(helpMessage)
			os.Exit(0)
		case "-e", "--execute":
			execute = flagValue
			args = args[2:]
		case "-c", "--content":
			contentFile = flagValue
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

	readDataFile(dataFile, &bytes)

	var data = map[string]interface{}{}

	parseDataFile(bytes, dataFormat, data)

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

	readContentFile(contentFile, data)

	if execute == "" {
		tmpl := template.New("")

		for _, templateFile := range templateFiles {
			templateSource, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.Fatal(err)
			}

			tmpl.Parse(string(templateSource))
		}

		tmpl.Execute(output, data)
	} else {
		tmpl, err := template.New("").ParseFiles(templateFiles...)
		if err != nil {
			log.Fatal(err)
		}

		tmpl.ExecuteTemplate(output, execute, data)
	}
}

func readContentFile(contentFile string, data map[string]interface{}) {
	if contentFile != "" {
		bytes, err := ioutil.ReadFile(contentFile)
		if err != nil {
			log.Fatal(err)
		}
		data["Content"] = string(bytes)
		data["content"] = string(bytes)
	}
}

func readDataFile(dataFile string, bytes *[]byte) {
	var err error

	if dataFile == "" {
		*bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		*bytes, err = ioutil.ReadFile(dataFile)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func parseDataFile(bytes []byte, dataFormat string, data map[string]interface{}) {
	var err error

	switch dataFormat {
	case "json":
		err = json.Unmarshal(bytes, &data)
	case "yaml":
		err = yaml.Unmarshal(bytes, &data)
	default:
		log.Fatal(fmt.Errorf(
			`The data format "%s" is invalid or not supported, use "json" or "yaml"`,
			dataFormat,
		))
	}
	if err != nil {
		log.Fatal(err)
	}
}
