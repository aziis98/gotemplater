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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var helpMessage = strings.TrimLeftFunc(`
gotemplater - A super small CLI utility wrapping template.Execute()
              and template.ExecuteTemplate()

USAGE: 
gotemplater [options...] [template files...]

EXAMPLES:

    Print to stdout: 
        gotemplater -f yaml -d data.yaml template.html

    Print to file: 
        gotemplater -d data.json -o rendered.html template.html

    Get data from file in another format: 
        gotemplater -f yaml -d data.yaml -o rendered.html template.html

    Get content from stdin, data from config.json and execute "main" defined in the given templates: 
        gotemplater -d config.json -c - -e main template-1.html template-2.html > rendered.html

OPTIONS:
    -h, --help			Shows this message

    -o, --output <file>		File to write to, by default uses stdout

    -e, --execute <name>	Template name to pass to .ExecuteTemplate(), otherwise uses .Execute()
    -c, --content <file>	Adds a content and Content variable to the context of the template 
				for rendering importing data from files, can be - for stdin


    -d, --data <file>		Data file to use
    -f, --format <format>	Format for the data file, can be json or yaml. By default it's json.

`, unicode.IsSpace)

func main() {
	args := os.Args[1:]
	templateFiles := []string{}

	templateName := ""
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
			templateName = flagValue
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

	data := getData(dataFile, dataFormat)
	output := getOutputWriter(outputFile)

	content := getContentFile(contentFile)
	if len(content) > 0 {
		data["Content"] = template.HTML(content)
		data["content"] = template.HTML(content)
	}

	executeTemplates(templateName, templateFiles, output, data)
}

func getOutputWriter(file string) io.Writer {
	if file == "" {
		return os.Stdout
	}

	w, err := os.Create(file)
	check(err)

	return w
}

func executeTemplates(templateName string, templateFiles []string, w io.Writer, data interface{}) {
	if templateName == "" {
		tmpl := template.New("")

		for _, templateFile := range templateFiles {
			templateSource, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.Fatal(err)
			}

			tmpl.Parse(string(templateSource))
		}

		tmpl.Execute(w, data)
	} else {
		tmpl, err := template.New("").ParseFiles(templateFiles...)
		check(err)

		tmpl.ExecuteTemplate(w, templateName, data)
	}
}

func getContentFile(file string) string {
	switch file {
	case "":
		return ""
	case "-":
		bytes, err := ioutil.ReadAll(os.Stdin)
		check(err)

		return string(bytes)
	default:
		bytes, err := ioutil.ReadFile(file)
		check(err)

		return string(bytes)
	}
}

func getData(file string, format string) (data map[string]interface{}) {
	var bytes []byte
	data = make(map[string]interface{})

	if file != "" {
		var err error

		bytes, err = ioutil.ReadFile(file)
		check(err)
	}
	if bytes == nil {
		return
	}

	switch format {
	case "json":
		err := json.Unmarshal(bytes, &data)
		check(err)
	case "yaml":
		err := yaml.Unmarshal(bytes, &data)
		check(err)
	default:
		log.Fatal(fmt.Errorf(
			`The data format "%s" is invalid or not supported, use "json" or "yaml"`,
			format,
		))
	}

	return
}
