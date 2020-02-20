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
	// Title is the title of the category
	// Link is the link of the category
	// Products is an array of all products in the category
	Title    string    `json:"categoryTitle"`
	Link     string    `json:"categoryLink"`
	Products []product `json:"categoryProducts"` 
}

type product struct {
	// Name is the name of the product
	// Image is the image link associated with the product
	// Price is the price of the product
	// Link is the link to the product
	Name  string `json:"productName"`
	Image string `json:"productImage"`
	Price string `json:"productPrice"`
	Link string `json:"productLink"`
}

type jSONData struct {
	// costcoProducts is a list of all products available from Costco
	CostcoProducts []categories `json:"allProducts"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector and echo object
	e := echo.New()
	c := colly.NewCollector(colly.Async(true))

	c.Limit(&colly.LimitRule{
		DomainGlob: "costco.com/*",
		RandomDelay: 2 * time.Second,
		Parallelism: 2,
	})

	var ls jSONData

	// Get the Category
	categorySelector := "#search-results > div.c_379317 > div"
	var datalist []categories

	c.OnHTML(categorySelector, func(e *colly.HTMLElement) {
		fmt.Println("First html before for each")
		e.ForEach("#search-results > div.c_379317 > div > div > div > div", func(_ int, h *colly.HTMLElement) {
			var products []product

			categoryName := e.ChildText("#search-results > div.c_379317 > div > div > div > div > a > div.h5-style-guide.eco-ftr-6across-title")
			categoryLink := e.ChildAttr("div > a", "href")
			fmt.Println(categoryName,"\n",categoryLink)
			fmt.Println("One HTML")
			// var categoryCount *int
			// categoryCount += 1

			d := categories{Title: categoryName, Link: categoryLink, Products: products}
			datalist = append(datalist, d)

			h.ForEach("#search-results > ctl:cache > div.product-list.grid", func(_ int, g *colly.HTMLElement) {
				productName := e.ChildText("#search-results > ctl:cache > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > p.description > a")
				productPrice := e.ChildText("#search-results > ctl:cache > div.product-list.grid > div > div > div.thumbnail > div.caption.link-behavior > div.caption > div > div")
				fmt.Println("Second HTML")

				p := product{Name: productName, Price: productPrice}
				datalist = append(datalist.Products, p)
			})
					
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping Costco
	c.Visit("https://www.costco.com/all-costco-grocery.html")

	ls = jSONData{CostcoProducts: datalist}
	fmt.Println(ls)

	// Serve to echo
	e.GET("/scrape", func(f echo.Context) error {
		return f.JSON(http.StatusOK, ls)
	})

	// Handle errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// After data is scraped, marshall to JSON
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

	c.Wait()
}
