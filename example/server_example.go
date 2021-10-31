package main

import (
	"github.com/panzg123/mouserpc"
)

func main() {
	s := mouserpc.NewServer()
	s.ListenAndServer()
}
