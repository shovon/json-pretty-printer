# JSON Pretty-Printer in Golang

Take in JSON, and output HTML-formated JSON, to get:

- syntax highlighted JSON
- pretty-formatted JSON, so that poorly formatted JSON code (with spaces and indentation missing) is now easier to read!

Developed for my CMPT 383 course at SFU.

## Usage

```
go build -o pretty-printer
./pretty-printer <path/to/file.json> > <path/to/output.html>
```

Then open your HTML file in a browser, and you should be good to go.

## License

```
Copyright © 2017 Salehen Shovon Rahman <salehen.rahman@gmail.com>
```
