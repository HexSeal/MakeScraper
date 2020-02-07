package main

import (
	"fmt"
	// "os"

	"github.com/gocolly/colly"
)

type post struct {
	Title string
	Content string
	LinkedPost string
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	// Get the post
	c.OnHTML("#thing_t3_ezyl5r", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		TopControversialPost := post{Title: e.Text}
		fmt.Println("\nOriginal Post:")
		fmt.Printf(link, "\n", TopControversialPost, "\n")
	})

	// Get the link of the post
	c.OnHTML("#thing_t3_ezyl5r > div.entry.unvoted > div > p.title > a", func(e *colly.HTMLElement) {
				link := e.Attr("href")

		// Print link
				fmt.Println("\nLinked Source:")
				fmt.Printf("%q -> %s\n", e.Text, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping reddit
	c.Visit("https://old.reddit.com/r/politics/controversial/")

		// Save the title
		// TopControversialPost := post{Title, e.Text}
		// os.Stdout
}
