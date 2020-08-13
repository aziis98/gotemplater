
# Go Templater

[One file CLI utility](./main.go) to render Go templates using _"data" files_ (in Json or Yaml) and _text content_ from files or stdin.

![Imgur](https://i.imgur.com/gkeFCO0.png)
<p align="center">
    An example of piping content down a pipeline to
    <br>
    build a static web page, the shell script <a href="./example/chaining/build.sh">is here</a>
</p>

**Why?.** I needed an extremely simple way to render templates with a CLI as I am moving my sites to use a Makefile instead of a full static site generator and I'd like a lot of flexibilty.

Under here there is a small documentation but it is probabily faster if you read the code directly as it is about 200 lines.

## Usage

    gotemplater [options...] [template files...]

### Options

- `-h`, `--help`
    Shows this message
- `-o`, `--output <file>`
    File to write to, by default uses stdout
- `-e`, `--execute <name>`
    Template name to pass to .ExecuteTemplate(), otherwise uses .Execute()
- `-c`, `--content <file>`
    Adds a `content` and `Content` variable to the context of the template for rendering importing data from files, can be `-` for stdin
- `-d`, `--data <file>`
    Data file to use
- `-f`, `--format <format>`
    Format for the data file, can be `json` or `yaml`. By default it's `json`.

### Examples

- Print to stdout: 
    
        gotemplater -d data.json template.html 

- Print to file: 
    
        gotemplater -d data.json -o rendered.html template.html

- Get the content from stdin: 
    
        gotemplater -f yaml -d data.yaml -c - template.html > rendered.html

- Example of chaining in the image on top: [chaining](./example/chaining/)


## TODO

The goal is to keep the project extremely small but if some cool ideas come up I will consider adding them. 

- Improve the names of the various functions
- Extract flags to a struct
- Maybe support TOML and some other formats
