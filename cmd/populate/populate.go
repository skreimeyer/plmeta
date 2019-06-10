// populate scrapes reads command line arguments and searches rosettacode.org
// for blocks of code listed under headings (programming language names) that
// match those strings. The code blocks are saved under individual files with
// extensions matched to the language or ".txt" as a default if the extension
// is unknown.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func main() {
	index := "http://rosettacode.org/wiki/Category:Programming_Tasks"
	langs := os.Args[1:]
	// normalize language names
	for i, l := range langs {
		langs[i] = slugUp(l)
		folder := filepath.Join("sources", langs[i])
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
	}
	log.Info("Starting...")

	// scrapes programming challenges
	taskCollector := colly.NewCollector(
		colly.AllowedDomains("rosettacode.org"),
		colly.MaxDepth(1),
	)
	// scrapes code blocks from a specific challenge
	codeCollector := colly.NewCollector(
		colly.AllowedDomains("rosettacode.org"),
		colly.Async(true),
		colly.MaxDepth(0),
	)
	// get tasks
	taskCollector.OnHTML("div.mw-category", func(e *colly.HTMLElement) {
		task := e.ChildAttrs("a", "href")
		log.Info("tasks found:", len(task))
		for i, t := range task {
			codeCollector.Visit(e.Request.AbsoluteURL(t))
			if i > 0 {
				break
			}
		}
	})
	// get code blocks
	codeCollector.OnRequest(func(r *colly.Request) {
		log.Info("Visiting:\t", r.URL.String())
	})

	codeCollector.OnResponse(func(r *colly.Response) {
		log.Info("Received response:\t", r.StatusCode)
	})

	codeCollector.OnHTML("h2", func(e *colly.HTMLElement) {
		i, ok := find(strings.ToUpper(e.ChildAttr("span", "id")), langs)
		if ok {
			log.Info("Found a match")
			// Get all pre.highlighted_source blocks up to the next title block
			sel := e.DOM.NextUntil("h2")
			blocks := sel.Filter("pre.highlighted_source")
			// assign common variables for writing
			algoURL := strings.Split(e.Request.URL.Path, "/")
			algo := algoURL[len(algoURL)-1]
			ext, ok := findExt(langs[i])
			if !ok {
				ext = "txt"
			}
			// We need to preserve <br> tags to make newlines.
			log.Info("Writing to files...")
			blocks.Each(func(j int, s *goquery.Selection) {
				fileName := fmt.Sprint(algo, "_", j, ".", ext)
				path := filepath.Join("sources", langs[i], fileName)
				f, err := os.Create(path)
				defer f.Close()
				if err != nil {
					log.Error(err)
					return
				}
				_, err = f.WriteString(customText(s))
				if err != nil {
					log.Error(err)
				}
			})
			//

		}
	})
	log.Info("visiting: ", index)
	taskCollector.Visit(index)

}

// get the index of a matched string within a slice of strings
func find(s string, col []string) (n int, ok bool) {
	for i, c := range col {
		if s == c {
			return i, true
		}
	}
	return 0, false
}

// slugify
func slugUp(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}

// container for map to match languages with their extensions
func findExt(lang string) (e string, ok bool) {
	ext := map[string]string{
		"PYTHON": "py",
		"PERL":   "pl",
		"PERL_6": "p6",
		"RUBY":   "rb",
		"PHP":    "php",
		"C":      "c",
		"C++":    "cpp",
		"JAVA":   "java",
		"SCHEME": "scm",
		"RUST":   "rs",
		"GO":     "go",
	}
	e, ok = ext[lang]
	return e, ok
}

// We need a customized modification of the Text method from goquery.
// This will preserve <br> tags as newline characters to keep our output
// code valid.
func customText(s *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// We will get weird characters from the nbsp that we need to get rid
			// of. We may get additional white space, but that's not such a problem.
			d := strings.ReplaceAll(n.Data, "\xa0", "\x20")
			d = strings.ReplaceAll(d, "\xc2", "\x20")
			buf.WriteString(d)
		}
		// insert a newline where there is a <br> tag
		if n.Data == "br" {
			buf.WriteString("\n")
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return buf.String()
}
