package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	Upper     bool
	Lower     bool
	Trim      bool
	BlockSize int
}

type Spaces struct {
	char rune
	amount int
}

func ParseFlags() (*Options, error) {
	var opts Options
	var convey string
	const optimalSize = 100

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset to start reading")
	flag.IntVar(&opts.Limit, "limit", 0, "limit to read")
	flag.StringVar(&convey, "conv", "", "text convertion")
	flag.IntVar(&opts.BlockSize, "block-size", optimalSize, "block for reading/writing")

	flag.Parse()

	//check if source is exists
	if len(opts.From) > 0 {
		if _, err := os.Stat(opts.From); errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	//check if destination is exists
	if len(opts.To) > 0 {
		if _, err := os.Stat(opts.To); !errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("destination file exists")
		}
	}
	//check if offset is valid
	if opts.Offset < 0 {
		return nil, errors.New("offset is negative")
	}
	if len(opts.From) > 0 {
		if f, err := os.Stat(opts.From); err != nil && f.Size() > int64(opts.Offset) {
			return nil, errors.New("offset is bigger than file")
		}
	}
	//check if limit is positive
	if opts.Limit < 0 {
		return nil, errors.New("limit is negative")
	}

	//check if block-size if positive
	if opts.BlockSize < 0 {
		return nil, errors.New("block-size is negative")
	}

	//check if convey is possible
	opts.Upper = false
	opts.Lower = false
	opts.Trim = false
	if len(convey) > 0 {
		var options []string = strings.Split(convey, ",")
		for _, option := range options {
			switch option {
			case "upper_case":
				opts.Upper = true
			case "lower_case":
				opts.Lower = true
			case "trim_spaces":
				opts.Trim = true
			default:
				return nil, errors.New("unknown option in conv: " + option)
			}
		}
		if opts.Lower && opts.Upper {
			return nil, errors.New("lower and upper options at the same time")
		}
	}
	return &opts, nil
}

var atStart bool = true

func DataDefinition(opts *Options) error {
	//create normal reader from file or input
	var reader io.Reader
	var readingFile *os.File
	var err error
	if len(opts.From) == 0 {
		reader = os.Stdin
	} else {
		readingFile, err = os.Open(opts.From)
		if err != nil {
			return err
		}
		reader = readingFile
		defer readingFile.Close()
	}
	if reader == nil {
		return errors.New("Cannot initialize Reader")
	}

	//adjust if offset is set
	if opts.Offset > 0 {
		var buf []byte = make([]byte, opts.Offset)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return errors.New("offset is to big")
			}
			return err
		}
	}

	//adjust if limit is set
	if opts.Limit != 0 {
		reader = io.LimitReader(reader, int64(opts.Limit))
		if reader == nil {
			return errors.New("Cannot initialize LimitReader")
		}
	}

	var writer io.Writer
	var writingFile *os.File
	if len(opts.To) == 0 {
		writer = os.Stdout
	} else {
		writingFile, err = os.OpenFile(opts.To, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		writer = writingFile
		defer writingFile.Close()
	}
	if writer == nil {
		return errors.New("Cannot initialize Writer")
	}

	var uncomplete [] byte
	var spaces [] Spaces = make([]Spaces, 0)
	
	for {
		var block []byte = make([]byte, opts.BlockSize)
		readSize, err := io.ReadFull(reader, block)
		if (readSize != len(block)) {
			block = block[:readSize]
		}
		uncomplete, block, spaces = adjustBuffer(opts, append(uncomplete, block...), spaces)

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				writer.Write(block)
				if uncomplete != nil {
					writer.Write(uncomplete)
				}
				break
			}
			return err
		}
		writer.Write(block)
	}
	return nil
}

func adjustBuffer (opts * Options, buf [] byte, prev [] Spaces) ([]byte, []byte, []Spaces) {
	if len(buf) == 0 {
		return nil, []byte{}, nil
	}
	var result []byte = make([]byte, 0)
	var str string = ""
	
	for len(buf) > 0 {
		r, size := utf8.DecodeRune(buf)
		if (r != utf8.RuneError) {
			break
		}
		result = append(result, buf[:size]...)
		buf = buf[size:]
	}
	if len(buf) == 0 {
		//all buf is an inconsist rune
		return result, decodeSpaces(prev), nil
	}
	
	r, size := utf8.DecodeRune(buf)
	for ; r != utf8.RuneError && buf[0] != '\x00'; r, size = utf8.DecodeRune(buf) {
		var s string = string(r)
		if (opts.Upper) {
			s = strings.ToUpper(s)
		} else if (opts.Lower) {
			s = strings.ToLower(s)
		}
		str += s
		buf = buf[size:]
	}
	
	var incomplete []byte
	if (size == 1 && r == utf8.RuneError) {
		//incomplete rune
		incomplete = []byte(strings.TrimRight(string(buf), "\x00"))
	} else {
		incomplete = nil
	}
	if (!opts.Trim) {
		return incomplete, append(result, []byte(str)...), nil
	}
	st, content, end := parseString(str)
	newSpaces := decodeSpaces(append(prev, st...))
	result = append(result, []byte(content)...)
	if content != "" {
		//write all spaces
		if (atStart) {
			atStart = false
			return incomplete, result, end
		}
		return incomplete, append(newSpaces, result...), end 
	}
	return incomplete, []byte{}, append(prev, st...)
}

func decodeSpaces (spaces [] Spaces) ([] byte) {
	buf := make([]byte, 0)
	for _, sp := range spaces {
		for sp.amount > 0 {
			buf = append(buf, byte(sp.char))
			sp.amount--
		}
	}
	return buf
}

func parseString(str string) ([] Spaces, string, [] Spaces) {
	startingSpaces := make([]Spaces, 0)
	var i int = 0
	for r, size := utf8.DecodeRune([]byte(str)); i < len(str) && unicode.IsSpace(r); r, size = utf8.DecodeRune([]byte(str)) {
		if len(startingSpaces) > 0 && startingSpaces[len(startingSpaces) - 1].char == r {
			startingSpaces[len(startingSpaces) - 1].amount++
		} else {
			startingSpaces = append(startingSpaces, Spaces{r, 1})
		}
		str = str[size:] 
	}
	if i == len(str) {
		return startingSpaces, "", make([]Spaces, 0)
	}

	endingSpaces := make([] Spaces, 0)
	i = len(str) - 1
	for r, size := utf8.DecodeLastRune([]byte(str)); unicode.IsSpace(r); r, size = utf8.DecodeLastRune([]byte(str)) {
		if len(endingSpaces) > 0 && endingSpaces[len(endingSpaces) - 1].char == r {
			endingSpaces[len(endingSpaces) - 1].amount++
		} else {
			endingSpaces = append(endingSpaces, Spaces{r, 1})
		}
		str = str[:len(str) - size]
	}
	return startingSpaces, strings.TrimSpace(str), endingSpaces
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = DataDefinition(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
