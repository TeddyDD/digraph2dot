package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/emicklei/dot"
)

var flagAttrs = flag.String("attr", "", "node attributes in format k=v[,k2=v2...]")

// from digrap source code
func split(line string) ([]string, error) {
	var (
		words   []string
		inWord  bool
		current bytes.Buffer
	)

	for len(line) > 0 {
		r, size := utf8.DecodeRuneInString(line)
		if unicode.IsSpace(r) {
			if inWord {
				words = append(words, current.String())
				current.Reset()
				inWord = false
			}
		} else if r == '"' {
			var ok bool
			size, ok = quotedLength(line)
			if !ok {
				return nil, errors.New("invalid quotation")
			}
			s, err := strconv.Unquote(line[:size])
			if err != nil {
				return nil, err
			}
			current.WriteString(s)
			inWord = true
		} else {
			current.WriteRune(r)
			inWord = true
		}
		line = line[size:]
	}
	if inWord {
		words = append(words, current.String())
	}
	return words, nil
}

// from digrap source code
func quotedLength(input string) (n int, ok bool) {
	var offset int

	// next returns the rune at offset, or -1 on EOF.
	// offset advances to just after that rune.
	next := func() rune {
		if offset < len(input) {
			r, size := utf8.DecodeRuneInString(input[offset:])
			offset += size
			return r
		}
		return -1
	}

	if next() != '"' {
		return // error: not a quotation
	}

	for {
		r := next()
		if r == '\n' || r < 0 {
			return // error: string literal not terminated
		}
		if r == '"' {
			return offset, true // success
		}
		if r == '\\' {
			var skip int
			switch next() {
			case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
				skip = 0
			case '0', '1', '2', '3', '4', '5', '6', '7':
				skip = 2
			case 'x':
				skip = 2
			case 'u':
				skip = 4
			case 'U':
				skip = 8
			default:
				return // error: invalid escape
			}

			for i := 0; i < skip; i++ {
				next()
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stderr)
	s := bufio.NewScanner(os.Stdin)
	graph := dot.NewGraph(dot.Directed)

	attrs := make(map[string]string)
	for _, kv := range strings.Split(*flagAttrs, ",") {
		if k, v, ok := strings.Cut(kv, "="); ok {
			attrs[k] = v
		}
	}

	setAttrs := func(n *dot.Node) {
		for k, v := range attrs {
			n.Attr(k, v)
		}
	}

	addEdges := func(from string, to ...string) {
		f := graph.Node(from)
		setAttrs(&f)
		for _, to := range to {
			t := graph.Node(to)
			setAttrs(&t)
			graph.Edge(f, t)
		}
	}

	var linenum int
	for s.Scan() {
		linenum++
		line := s.Text()
		nds, err := split(line)
		if err != nil {
			log.Fatalf("line %d: %s", linenum, err.Error())
		}
		if len(nds) > 0 {
			addEdges(nds[0], nds[1:]...)
		}
	}
	graph.Write(os.Stdout)
}
