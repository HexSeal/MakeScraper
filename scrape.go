package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
)

type categories struct {
		Title	string	`json:"categoryTitle"`
		Link	string	`json:"categoryLink"`
		Products	[]product	`json:"categoryProducts`
}

type product struct {
	Name	string	`json:"productName"`
	Image	string	`json:"productImage"`
	Price	string	`json:"productPrice"`
}

type JSONData struct {
	costcoProducts []categories `json:"allProducts"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector and echo object
	e := echo.New()
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Delay: 1 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	// var datalist []data
	// var b data

	// #search-results > div.c_379317 > div > div > div > div

	// Get the Category
	categorySelector := "#search-results > div.c_379317 > div"
	var datalist []categories
	var d categories

	c.OnHTML(categorySelector, func(e *colly.HTMLElement) {
		e.ForEach("#search-results > div.c_379317 > div > div > div > div", func(_ int, h *colly.HTMLElement) {
			var products []product

			categoryName := e.ChildText("#search-results > div.c_379317 > div > div > div > div > a > div.h5-style-guide.eco-ftr-6across-title")
			categoryLink := e.ChildAttr("div > a", "href")
			// var categoryCount *int
			// categoryCount += 1

			// Here I want to go through each link and scrape all products

			d = categories{Title: categoryName, Link: categoryLink, Products: products}
			datalist = append(datalist, d)
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping Costco
	c.Visit("https://www.costco.com/all-costco-grocery.html")

	ls := JSONData{costcoProducts: datalist}

	e.GET("/scrape", func(f echo.Context) error {
		return f.JSON(http.StatusOK, ls)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	DataJSONarr, err := json.MarshalIndent(ls, "", "	")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("output.json", DataJSONarr, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello")
	e.Logger.Fatal(e.Start(":8000"))
}

