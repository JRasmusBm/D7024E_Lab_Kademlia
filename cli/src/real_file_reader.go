package main

import (
	"io/ioutil"
)

type RealFileReader struct {
}

func (r *RealFileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
