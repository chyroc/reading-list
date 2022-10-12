package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
const dateTimeFormat = "2006-01-02 15:04:05"

func urlHash(uri string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(uri)))
}

type readingListEntry struct {
	URL   string `csv:"url,omitempty"`
	Title string `csv:"title,omitempty"`
	Date  string `csv:"date,omitempty"`
}

func loadData() ([]*readingListEntry, error) {
	files, err := ioutil.ReadDir("./data")
	if err != nil {
		return nil, err
	}

	entries := make([]*readingListEntry, 0, len(files))
	for _, v := range files {
		if v.IsDir() || !strings.HasSuffix(v.Name(), ".json") {
			continue
		}
		bs, err := ioutil.ReadFile("./data/" + v.Name())
		if err != nil {
			return nil, err
		}
		var entry readingListEntry
		err = json.Unmarshal(bs, &entry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}
