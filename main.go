package main

import (
	"io"
	"os"
	"log"
	"flag"
	"bufio"
	"regexp"
	"strings"
)

var (
	splitter = regexp.MustCompile(`(\s+|\p{L}+|[^\p{L}]+|\S+)`)
)

var (
	inputFile = flag.String("input", "", "input file (default is stdin)")
	outputFile = flag.String("output", "", "output file (default is stdout)")
	dictFile = flag.String("dictionary", "", "dictionary file")
)

func main() {
	flag.Parse()

	r := os.Stdin
	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal("open input file:", err)
		}
		r = f
	}

	if *dictFile != "" {
		f, err := os.Open(*dictFile)
		if err != nil {
			log.Fatal("open dictionary:", err)
		}
		scanner := bufio.NewScanner(f)
		for line := 1; scanner.Scan(); line++ {
			b := scanner.Text()
			kv := strings.Split(b, ":")
			if len(kv) != 2 {
				log.Fatalf("invalid line %d: %q", line, kv)
			}
			dictionary[kv[0]] = kv[1]
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("reading dictionary:", err)
		}
	}

	b := &strings.Builder{}
	if _, err := io.Copy(b, r); err != nil {
		log.Fatal("read input file:", err)
	}
	text := b.String()

	w := os.Stdout
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal("create output file:", err)
		}
		w = f
	}

	tokens := splitter.FindAllString(text, -1)
	for _, tok := range tokens {
		if val, ok := dictionary[tok]; ok {
			tok = val
		}
		io.WriteString(w, tok)
	}
}
