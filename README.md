
# Gotemplater

[One file CLI utility](./main.go) to render Go templates to stdout and files from JSON and YAML "data" files and text content from files and stdin.

![Imgur](https://i.imgur.com/gkeFCO0.png)
<p align="center">
    An example of piping content down a pipeline to
    <br>
    build a static web page, the shell script <a href="./example/chaining/build.sh">is here</a>
</p>

**Why?.** I needed an extremely simple way to render templates with a CLI as I am moving my sites to use a Makefile instead of a full static site generator and I'd like a lot of flexibilty.

Under here there is a small documentation but it is probabily faster if you read the code directly as it is about 200 lines.

## Options

TODO

## TODO

The goal is to keep the project extremely small but if some cool ideas come up I will consider adding them. 

- Improve the names of the various functions
- Extract flags to a struct
- Maybe support TOML and some other formats