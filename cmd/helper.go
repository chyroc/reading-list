package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/go-playground/validator"
)

var validate = validator.New()

type Inputs struct {
	URL string `validate:"required,url" query:"url"`
}

func (i *Inputs) Validate() error {
	return validate.Struct(i)
}

const dateFormat = "2006-01-02"

func urlHash(uri string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(uri)))
}

type readingListEntry struct {
	URL   string `csv:"url,omitempty"`
	Title string `csv:"title,omitempty"`
	Date  string `csv:"date,omitempty"`
}
