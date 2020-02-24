package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
)

type categories struct {
	// Title is the title of the category
	// Link is the link of the category
	// Products is an array of all products in the category
	// Image is the thumbnail image of the category
	Title   string	`json:"categoryTitle"`
	Link    string	`json:"categoryLink"`
	Image	string	`json:"categoryImage"`
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

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector and echo object
	e := echo.New()
	c := colly.NewCollector(colly.Async(true))

	// Limitations so our 
	c.Limit(&colly.LimitRule{
		// DomainGlob: "https://www.bjs.com/*",
		RandomDelay: 2 * time.Second,
		Parallelism: 2,
	})

	// Get the Category
	categorySelector := "#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div"
	var categoriesList []categories
	var p []product

	var categoryName string
	var categoryLink string
	var categoryImage string

	// Scrape categories
	c.OnHTML(categorySelector, func(e *colly.HTMLElement) {
		e.ForEach("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div", func(num int, h *colly.HTMLElement) {
			// var p []product

			// fmt.Println("Category Number:", num)

			categoryName = e.ChildText("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div:nth-child(" + strconv.Itoa(num + 1) + ") > a")
			categoryLink = e.ChildAttr("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div:nth-child(" + strconv.Itoa(num + 1) + ") > a", "href")
			categoryImage = e.ChildAttr("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div:nth-child("+ strconv.Itoa(num + 1) + ") > a > img", "src")
			// fmt.Println(categoryName, "\n", categoryLink)

			// Append the category data to the struct
			d := categories{Title: categoryName, Link: categoryLink, Image: categoryImage, Products: p}
			categoriesList = append(categoriesList, d)

			// Get the link for each category and visit it with the next OnHTML request so we can scrape all the products of said category
			e.Request.Visit(categoryLink) // Add &pagesize=80 to get max number of products per page
		})
		//fmt.Println(categoriesList)
	})

	// fmt.Println("product onhtml start")

	// For each category, scrape all product data by following the category link
	c.OnHTML("#contentOverlay > div > app-cat-plp-page > div > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div" , func(b *colly.HTMLElement) {
		fmt.Println("Product OnHTML request goes off")
		b.ForEach("#contentOverlay > div > app-cat-plp-page > div > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div > div", func(count int, g *colly.HTMLElement) {
			fmt.Println(count)
			productName := g.ChildText("#contentOverlay > div > app-cat-plp-page > div > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div > div > app-product-card > div > a.product-link > h2.product-title.section.d-none.d-sm-block")
			productPrice := g.ChildText("#contentOverlay > div > app-cat-plp-page > div > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div > div > app-product-card > div > div.price-block.section > div.display-price > span")
			productImage := g.ChildAttr("#contentOverlay > div > app-cat-plp-page > div > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div > div > app-product-card > div > a.section.img-link > img", "src")
			productLink := g.ChildAttr("#contentOverlay > div > app-cat-plp-page > div:nth-child(1) > app-search-result-page-gb > div.bottomContainer > div.rightSection.show-mobile > div.rightBottom > app-products-container > div > div > div > app-product-card > div > a.product-link", "href")

			// Adding individual products to the product list
			pl := product{Name: productName, Price: productPrice, Image: productImage, Link: productLink}
			fmt.Println(pl)
			p = append(p, pl)
		})	
	})


	// Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// })

	// Start scraping BJ's wholesale site
	c.Visit("https://www.bjs.com/content?template=B&espot_main=EverydayEssentials&source=megamenu")

	// Serve to echo
	e.GET("/scrape", func(f echo.Context) error {
		return f.JSON(http.StatusOK, categoriesList)
	})

	// Handle errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Alert when done scraping
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// After data is scraped, marshall to JSON
	DataJSONarr, err := json.MarshalIndent(categoriesList, "", "	")
	if err != nil {
		panic(err)
	}

	// Writing the marshalled JSON data to output.json
	err = ioutil.WriteFile("output.json", DataJSONarr, 0644)
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8000"))

	c.Wait()
}
