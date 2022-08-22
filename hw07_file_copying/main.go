package main

import (
	"flag"
)

var (
	from, to      string
	limit, offset int64
	isAsync       bool
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.BoolVar(&isAsync, "async", false, "is async mode")
}

func main() {
	flag.Parse()

	err := Copy(from, to, offset, limit, isAsync)
	if err != nil {
		errorLog.Fatal(err)
	}
}
