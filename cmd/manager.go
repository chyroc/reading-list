package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/reading-list/transport"
	"github.com/jszwec/csvutil"
)

func AddRowToCSV() error {
	data := new(transport.Inputs)
	if err := loadInputs(data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return err
	}

	title, err := getHtmlTitle(data.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get html title for URL %s\n", data.URL)
	}

	// if CSV file does not exist, create it with a header
	csvFilePath := readingListFile
	{
		_, err := os.Stat(csvFilePath)
		if os.IsNotExist(err) {

			b, err := csvutil.Marshal([]*readingListEntry{})
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(csvFilePath, b, 0644)
			if err != nil {
				return err
			}

		}
	}

	// make changes to CSV file

	f, err := os.OpenFile(csvFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	b, err := csvutil.Marshal([]*readingListEntry{{
		URL:   data.URL,
		Title: title,
		Date:  time.Now(),
	}})
	if err != nil {
		return err
	}

	splitOutput := bytes.Split(b, []byte("\n"))

	if _, err = f.Write(append(splitOutput[1], byte('\n'))); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

func loadInputs(data any) error {
	return json.Unmarshal([]byte(os.Getenv("RL_INPUT_JSON")), data)
}

func getHtmlTitle(url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("getHtmlTitle: non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	return doc.Find("title").Text(), nil
}

var httpClient = &http.Client{
	Timeout: time.Second * 20,
}
