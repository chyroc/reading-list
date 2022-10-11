package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func AddURL() error {
	data := new(Inputs)
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

	name := fmt.Sprintf("./data/%s.json", urlHash(data.URL))
	body, _ := json.MarshalIndent(readingListEntry{
		URL:   data.URL,
		Title: title,
		Date:  time.Now().Format(dateFormat),
	}, "", "  ")

	f, _ := os.Stat(name)
	if f != nil {
		fmt.Println("url already exists")
		return nil
	}

	return ioutil.WriteFile(name, body, 0644)
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
