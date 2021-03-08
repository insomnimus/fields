package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/insomnimus/fields/args"
)

var (
	exeName = "fields"
	lines   []int
	words   []int

	separator        = ""
	maxLine          = -1
	file, word, line string
)

var (
	reSep  = regexp.MustCompile(`^\-(?:s|\-separator)=(.+)$`)
	reWord = regexp.MustCompile(`^\-(?:w|\-word)=(.+)$`)
	reLine = regexp.MustCompile(`^\-(?:l|\-line)=(.+)$`)
	reFile = regexp.MustCompile(`^\-(?:f|\-file)=(.+)$`)
)

func readStream(in io.Reader) {
	scanner := bufio.NewScanner(in)
	i := 1
	for scanner.Scan() {
		if maxLine > 0 && i > maxLine {
			break
		}
		evaluate(scanner.Text(), i)
		i++
	}
}

func evaluate(s string, ln int) {
	if s == "" {
		return
	}
	if len(lines) > 0 {
		yes := false
		for _, l := range lines {
			if l > ln {
				return
			}
			if l == ln {
				yes = true
				break
			}
		}
		if !yes {
			return
		}
	}
	if len(words) == 0 {
		fmt.Println(s)
		return
	}
	var fields []string
	if separator == "" {
		fields = strings.Fields(s)
	} else {
		fields = strings.Split(s, separator)
	}
	var matches []string
	for _, w := range words {
		if w-1 >= len(fields) {
			break
		}
		matches = append(matches, fields[w-1])
	}
	if len(matches) == 0 {
		return
	}
	if len(matches) == 1 {
		fmt.Println(matches[0])
		return
	}
	fmt.Println(strings.Join(matches, "  "))
}

func showHelp() {
	log.Printf(`%s, select substrings and lines from a stream
usage:
	%s [options] arguments

options are:
	-w, --word=<numbers>: word indices to be printed, starts from 1
	-l, --line=<numbers>: line indices to be printed, starts from 1
	-s, --separator=<separator>: separator to split the lines by, default is spaces
	-f, --file=<file>: file to scan
	-h, --help: show this message

number ranges are accepted:
	--line=1..10,15,20..25
if the file is omitted, the input stream will be the standard input
if all the flags are omitted, argument order corresponds to:
	<words> <lines> <separator> <file>`,
		exeName, exeName)
	os.Exit(0)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("")
	exeName = filepath.Base(strings.TrimSuffix(os.Args[0], ".exe"))
	if len(os.Args) == 1 {
		showHelp()
	}
	arguments := os.Args[1:]
	var a string
LOOP:
	for i := 0; i < len(arguments); i++ {
		a = arguments[i]
		if a[0] == '-' {
			switch a {
			case "-s", "--separator":
				if i+1 >= len(arguments) {
					log.Fatal("error: the --separator flag was set but the value is not provided")
				}
				i++
				separator = arguments[i]
			case "-h", "--help":
				showHelp()
			case "-f", "--file":
				if i+1 >= len(arguments) {
					log.Fatal("error: the --file flag was set but the value is not provided")
				}
				i++
				file = arguments[i]
			case "-l", "--line":
				if i+1 >= len(arguments) {
					log.Fatal("error: the --line flag was set but the value is not provided")
				}
				i++
				line = arguments[i]
			case "-w", "--word":
				if i+1 >= len(arguments) {
					log.Fatal("error: the --word flag was set but the value is not provided")
				}
				i++
				word = arguments[i]
			default:
				if m := reSep.FindStringSubmatch(a); len(m) == 2 {
					separator = m[1]
					continue LOOP
				}
				if m := reWord.FindStringSubmatch(a); len(m) == 2 {
					word = m[1]
					continue LOOP
				}
				if m := reLine.FindStringSubmatch(a); len(m) == 2 {
					line = m[1]
					continue LOOP
				}
				if m := reFile.FindStringSubmatch(a); len(m) == 2 {
					file = m[1]
					continue LOOP
				}
				log.Fatalf("unknown command line option '%s'", a)
			}
			continue
		}
		switch {
		case word == "":
			word = a
		case line == "":
			line = a
		case separator == "":
			separator = a
		case file == "":
			file = a
		default:
			log.Fatal("too many arguments")
		}
	}
	var err error
	if line != "" {
		lines, err = args.Parse(line)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		if len(lines) == 0 {
			log.Fatalf("invalid value --line=%q", line)
		}
		maxLine = lines[len(lines)-1]
	}

	if word != "" {
		words, err = args.Parse(word)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		if len(words) == 0 {
			log.Fatalf("invalid value --word=%q", word)
		}
	}
	if file == "" {
		readStream(os.Stdin)
		return
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	readStream(f)
}
