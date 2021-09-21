package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	schema, err := os.Open("/Users/dino/Code/git/autobahn/services/management-server/src/schema.graphql")
	if err != nil {
		panic(err)
	}

	h, err := New(schema)
	if err != nil {
		panic(err)
	}
	println("Lets go!")
	done, err := h.Next()
	if err != nil {
		panic(err)
	}

	for !done {
		println(h.Token())
		done, err = h.Next()
		if err != nil {
			panic(err)
		}
	}
	if done {
		println("done")
	}

}

func New(r io.Reader) (hulud, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanGraphQLToken)
	return hulud{scanner: scanner}, nil
}

type hulud struct {
	scanner *bufio.Scanner
}

func (h *hulud) Next() (done bool, err error) {
	ok := h.scanner.Scan()
	if !ok {
		if h.scanner.Err() != nil {
			return false, h.scanner.Err()
		}
		return true, nil
	}
	return false, nil
}

func (h *hulud) Token() string {
	return h.scanner.Text()
}
