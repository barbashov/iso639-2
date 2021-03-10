package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/barbashov/iso639-2"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultInput                    = "http://loc.gov/standards/iso639-2/ISO-639-2_utf-8.txt"
	HttpTimeout                     = 60 * time.Second
	UTF8BOM                         = "\uFEFF"
	InputFileSeparator              = "|"
	InputFileLineColumns            = 5
	InputFileLanguageNamesSeparator = "; "
)

func main() {
	lookup := map[string]iso639_2.Language{}

	inputFile := flag.String("i", DefaultInput,
		fmt.Sprintf("Path or URL to input file in pipe-separated loc.gov format (default %s)", DefaultInput))
	outfile := flag.String("o", "", "Output file (default - standard output)")
	flag.Parse()

	rd := getInput(*inputFile)

	scanner := bufio.NewScanner(rd)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		columns := strings.Split(line, InputFileSeparator)
		if len(columns) != InputFileLineColumns {
			log.Fatalf("Error reading input file at line %d: %d columns expected, %d found",
				lineNum, InputFileLineColumns, len(columns))
		}

		if strings.HasPrefix(columns[0], UTF8BOM) {
			// remove UTF-8 BOM
			columns[0] = columns[0][len(UTF8BOM):]
		}

		language := iso639_2.Language{
			Alpha3:  columns[0],
			Alpha2:  columns[2],
			English: strings.Split(columns[3], InputFileLanguageNamesSeparator),
			French:  strings.Split(columns[4], InputFileLanguageNamesSeparator),
		}

		if language.Alpha3 != "" {
			lookup[language.Alpha3] = language
		}

		if language.Alpha2 != "" {
			lookup[language.Alpha2] = language
		}
	}

	wr := os.Stdout
	if *outfile != "" {
		var err error
		wr, err = os.Create(*outfile)
		if err != nil {
			log.Fatalf("Can't create output file '%s': %v", *outfile, err)
		}
	}

	outputLookup(wr, lookup)
}

func getInput(uri string) io.Reader {
	parsedUrl, err := url.Parse(uri)
	if err != nil || parsedUrl.Scheme == "" {
		f, err := os.Open(uri)
		if err != nil {
			log.Fatalf("Can't open input file '%s': %v", uri, err)
		}
		return bufio.NewReader(f)
	}

	httpClient := &http.Client{
		Timeout: HttpTimeout,
	}

	r, err := httpClient.Get(uri)
	if err != nil {
		log.Fatalf("Can't download input file '%s': %v", uri, err)
	}
	defer r.Body.Close()

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response from '%s': %v", uri, err)
	}

	return bytes.NewReader(bs)
}

func outputLookup(w io.Writer, lookup map[string]iso639_2.Language) {
	lookupStr := fmt.Sprintf("%#v", lookup)
	replacer := strings.NewReplacer(
		"iso639_2.Language{Alpha3", "{Alpha3",
		"},", "},\n",
		"map[string]iso639_2.Language{", "map[string]Language{\n",
	)
	lookupStr = replacer.Replace(lookupStr)

	buf := bytes.Buffer{}
	_, err := fmt.Fprintf(&buf, "package iso639_2\n\nvar Languages = %s\n", lookupStr)
	if err != nil {
		log.Fatalf("Error generating: %v", err)
	}

	outBytes, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Error formatting generated code: %v", err)
	}

	_, err = w.Write(outBytes)
	if err != nil {
		log.Fatalf("Error writing to output: %v", err)
	}
}
