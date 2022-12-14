package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func GenerateSite() error {
	entries, err := loadData()
	if err != nil {
		return err
	}

	numArticles := len(entries)
	groupedEntries := groupEntriesByMonth(entries)

	head := Div(
		H1(g.Text(pageTitle)),
		P(g.Raw(
			fmt.Sprintf(
				"A mostly complete list of articles I've read on the internet.<br>The articles on this list do not necessarily represent my views or opinions and should not be construed as doing so.<br><br>There are currently %d entries in the list<br>Last modified %s<br>Repo: %s",
				numArticles,
				time.Now().Format(dateFormat),
				"<a href=\"https://github.com/chyroc/reading-list\" rel=\"noopener\"><code>chyroc/reading-list</code></a>",
			),
		)),
	)

	listing := makeListHTML(groupedEntries)

	outputContent, err := renderHTMLPage(pageTitle, []g.Node{head, Hr(), listing})
	if err != nil {
		return err
	}

	const outputDir = ".site"
	_ = os.Mkdir(outputDir, 0777)
	err = ioutil.WriteFile(outputDir+"/index.html", outputContent, 0644)
	if err != nil {
		return err
	}

	rss, json_, err := generateRSS(entries)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outputDir+"/rss.xml", []byte(rss), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outputDir+"/rss.json", []byte(json_), 0644)
	if err != nil {
		return err
	}
	return nil
}

// renderHTMLPage renders a complete HTML page
func renderHTMLPage(title string, body []g.Node) ([]byte, error) {
	b := new(bytes.Buffer)
	err := c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en-GB",
		Head:     []g.Node{Link(g.Attr("rel", "stylesheet"), g.Attr("href", "https://www.tdpain.net/assets/css/ghpages.css"), g.Attr("type", "text/css"))},
		Body:     []g.Node{Div(g.Attr("class", "container"), g.Group(body))},
	}).Render(b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type entryGroup struct {
	Title   string
	Date    time.Time
	ID      string
	Entries entrySlice
}

type entryGroupSlice []*entryGroup

func (e entryGroupSlice) Len() int {
	return len(e)
}

func (e entryGroupSlice) Less(i, j int) bool {
	return e[i].Date.After(e[j].Date)
}

func (e entryGroupSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type entrySlice []*readingListEntry

func (e entrySlice) Len() int {
	return len(e)
}

func (e entrySlice) Less(i, j int) bool {
	ii, _ := time.Parse(dateTimeFormat, e[i].Date)
	jj, _ := time.Parse(dateTimeFormat, e[j].Date)
	return ii.After(jj)
}

func (e entrySlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func groupEntriesByMonth(entries []*readingListEntry) entryGroupSlice {
	groupMap := make(map[time.Time]*entryGroup)

	for _, entry := range entries {
		date, err := time.Parse(dateTimeFormat, entry.Date)
		if err != nil {
			continue
		}
		newTime := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
		if groupMap[newTime] == nil {
			groupMap[newTime] = &entryGroup{
				Date:  newTime,
				Title: fmt.Sprintf("%s %d", newTime.Month().String(), newTime.Year()),
				ID:    strings.ToLower(fmt.Sprintf("%s-%d", newTime.Month().String(), newTime.Year())),
			}
		}
		groupMap[newTime].Entries = append(groupMap[newTime].Entries, entry)
	}

	var o entryGroupSlice
	for _, group := range groupMap {
		sort.Sort(group.Entries)
		o = append(o, group)
	}
	sort.Sort(o)

	return o
}

// makeTILHTML generates HTML from a []*entryGroup to make a list of articles
func makeListHTML(groups []*entryGroup) g.Node {

	headerLevel := H3

	numGroups := len(groups)

	var subsections []g.Node
	for i := numGroups - 1; i >= 0; i -= 1 {
		group := groups[i]
		subsections = append(subsections, A(g.Attr("href", "#"+group.ID), g.Textf("%s %d", group.Date.Month().String()[:3], group.Date.Year())))
	}

	parts := []g.Node{
		Br(),
		Span(g.Text("Jump to :: "), g.Group(g.Map(len(subsections), func(i int) g.Node {
			n := subsections[i]
			if i != len(subsections)-1 {
				n = g.Group([]g.Node{n, g.Text(" :: ")})
			}
			return n
		}))),
	}

	for _, group := range groups {

		dateString := group.Title

		header := headerLevel(g.Attr("id", group.ID), g.Text(dateString))

		var entries []g.Node
		for _, article := range group.Entries {
			entries = append(entries, articleLinkComponent(
				article.URL,
				article.Title,
				article.Date),
			)
		}

		parts = append(parts, header, Ul(entries...))
	}

	return Div(parts...)
}

func articleLinkComponent(url, title, date string) g.Node {
	d, _ := time.Parse(dateTimeFormat, date)
	return Li(
		A(g.Attr("href", url), g.Text(title)),
		g.Text(" - "+d.Format(dateFormat)),
	)
}
