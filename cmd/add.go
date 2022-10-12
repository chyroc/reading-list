package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

	err1 := addOne(data.URL, time.Now().Format(dateFormat), false)
	err2 := reGenerate()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func addOne(url string, date string, ignoreData bool) error {
	fmt.Println("url", url)

	title, err := getHtmlTitle(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get html title for URL %s, %s\n", url, err)
	}
	fmt.Println("title", title)

	name := fmt.Sprintf("./data/%s.json", urlHash(url))
	body, _ := json.MarshalIndent(readingListEntry{
		URL:   url,
		Title: title,
		Date:  date,
	}, "", "  ")

	if !ignoreData {
		f, _ := os.Stat(name)
		if f != nil {
			fmt.Println("url already exists")
			return nil
		}
	}

	return ioutil.WriteFile(name, body, 0644)
}

func reGenerate() error {
	entities, err := loadData()
	if err != nil {
		return err
	}
	for _, v := range entities {
		if strings.TrimSpace(v.Title) == "" {
			fmt.Println("empty title", v.URL)
			err = addOne(v.URL, v.Date, true)
			if err != nil {
				return err
			}
		}
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

	title := strings.TrimSpace(doc.Find("title").Text())
	if title != "" {
		return title, nil
	}

	// og:title
	title = strings.TrimSpace(doc.Find("meta[property='og:title']").AttrOr("content", ""))
	if title != "" {
		return title, nil
	}

	// twitter:title
	title = strings.TrimSpace(doc.Find("meta[name='twitter:title']").AttrOr("content", ""))
	if title != "" {
		return title, nil
	}

	return "", nil
}

var httpClient = &http.Client{
	Timeout: time.Second * 20,
}
