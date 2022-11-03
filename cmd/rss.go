package main

import (
	"github.com/gorilla/feeds"
	"time"
)

func generateRSS(data []*readingListEntry) (string, string, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       pageTitle,
		Link:        &feeds.Link{Href: "https://reading-list.chyroc.cn/"},
		Description: pageTitle,
		Author:      &feeds.Author{Name: "chyroc"},
		Created:     now,
	}

	for _, v := range data {
		date, err := time.Parse(dateTimeFormat, v.Date)
		if err != nil {
			continue
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Title:     v.Title,
			Link:      &feeds.Link{Href: v.URL},
			Source:    nil,
			Id:        urlHash(v.URL),
			Updated:   date,
			Created:   date,
			Enclosure: nil,
			Content:   "",
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", "", err
	}
	json_, err := feed.ToJSON()
	if err != nil {
		return "", "", err
	}
	return rss, json_, nil
}
