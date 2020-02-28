package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	// "github.com/labstack/echo/v4"
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type categories struct {
	// Title is the title of the category
	// Link is the link of the category
	// Products is an array of all products in the category
	// Image is the thumbnail image of the category
	// gorm.Model
	Title    string    `json:"categoryTitle"`
	Link     string    `json:"categoryLink"`
	Image    string    `json:"categoryImage"`
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
	Link  string `json:"productLink"`
}

func scrape() {
	// Initialize gorm and the database
	// db, err := gorm.Open("sqlite3", "test.db")
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// defer db.Close()

	// Migrate the schema
	// db.AutoMigrate(&Category{})

	// Instantiate default collector
	c := colly.NewCollector()

	// Limitations so our bot can't accidentally dox someone.
	c.Limit(&colly.LimitRule{
		// DomainGlob: "https://www.bjs.com/*",
		RandomDelay: 1 * time.Second,
		// Parallelism: 2,
	})

	// Get the Category
	categorySelector := "#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div"
	var categoriesList []categories

	// Scrape categories
	c.OnHTML(categorySelector, func(e *colly.HTMLElement) {
		e.ForEach("#contentOverlay > div > app-content > div > div > div > div > div > div.bopic-hero > div > div > div > div > div > div", func(num int, h *colly.HTMLElement) {
			// fmt.Println("Category Number:", num)

			// Faster: create a string of common html elements, make those a variables and avoid memory-expensive concatenation
			categoryName := e.ChildText("div.bopic-hero > div > div > div > div > div > div:nth-child(" + strconv.Itoa(num+1) + ") > a")
			categoryLink := e.ChildAttr("div.bopic-hero > div > div > div > div > div > div:nth-child("+strconv.Itoa(num+1)+") > a", "href")
			categoryImage := e.ChildAttr("div.bopic-hero > div > div > div > div > div > div:nth-child("+strconv.Itoa(num+1)+") > a > img", "src")
			// fmt.Println(categoryName, "\n", categoryLink)

			//db.Create(&Category{Title: categoryName, Link: categoryLink, Image: categoryImage})
			//db.Create(&Category{Title: breakfast, Link: reddit.com, Image: https://i.kym-cdn.com/entries/icons/original/000/027/475/Screen_Shot_2018-10-25_at_11.02.15_AM.png})

			// Append the category data to the struct
			// productScrape := ScrapeProducts(categoryName, categoryLink)
			d := categories{Title: categoryName, Link: categoryLink, Image: categoryImage} // , Products: productScrape
			categoriesList = append(categoriesList, d)

			// e.Request.Visit(categoryLink)
		})
		fmt.Println(categoriesList)
	})

		// Start scraping bj's wholesale site
		c.Visit("https://www.bjs.com/content?template=B&espot_main=EverydayEssentials&source=megamenu")

		// Handle errors
		c.OnError(func(_ *colly.Response, err error) {
			log.Println("Something went wrong:", err)
		})

		// Alert when done scraping
		c.OnScraped(func(r *colly.Response) {
			//fmt.Println("Finished", r.Request.URL)
		})

		// After data is scraped, marshall to JSON
		DataJSONarr, err := json.MarshalIndent(categoriesList, "", "	")
		if err != nil {
			panic(err)
		}
		// fmt.Println(categoriesList)

		// Writing the marshalled JSON data to output.json
		err = ioutil.WriteFile("output.json", DataJSONarr, 0644)
		if err != nil {
			panic(err)
		}
		c.Wait()

		// GORM test, please ignore
		// // Test GORM by reading the entry with id 1
		// var testCat Category
		// db.First(&testCat, 1)
		// // Update - update product's price to 2000
		// db.Model(&testCat).Update("Price", 2000)
		// // Delete - delete product
		// db.Delete(&testCat)
}


// ScrapeProducts scrapes all product data by following the category link
func ScrapeProducts(categoryName string, categoryLink string) {
	// Instantiate default collector and echo object
	c := colly.NewCollector()

	// Limitations so our
	c.Limit(&colly.LimitRule{
		// DomainGlob: "https://www.bjs.com/*",
		RandomDelay: 1 * time.Second,
		// Parallelism: 2,
	})

	var categoryProducts []product
	// For each category, scrape all product data by following the category link
	c.OnHTML("app-products-container > div > div", func(b *colly.HTMLElement) {
		// fmt.Println("Product OnHTML request goes off")
		fmt.Printf("\nCurrent link: %s\n", b.Request.URL)
		b.ForEach("div.rightBottom > app-products-container > div > div > div", func(count int, g *colly.HTMLElement) {
			// fmt.Println(count)

			productName := g.ChildText("div:nth-child(" + strconv.Itoa(count+1) + ") > app-product-card > div > a.product-link > h2.product-title.section.d-none.d-sm-block")
			productLink := g.ChildAttr("div:nth-child("+strconv.Itoa(count+1)+") > app-product-card > div > a.product-link", "href")
			productPrice := g.ChildText("div:nth-child(" + strconv.Itoa(count+1) + ") > app-product-card > div > div.price-block.section > div.display-price > span")
			productImage := g.ChildAttr("div:nth-child("+strconv.Itoa(count+1)+") > app-product-card > div > a.section.img-link > img", "src")
			
			// Get the category for the product
			// Find the category inside categoriesList & save to theCategory
			// categoryName :=

			// theCategory := categoriesList[categoryName]
			// theCategoryLink := theCategory.Link

			// Adding individual products to the product list
			fmt.Printf("Product name: %v, Product Price: %v, Product Image: %v, Product Link: %v\n", productName, productPrice, productImage, productLink)
			productInstance := product{Name: productName, Price: productPrice, Image: productImage, Link: productLink}
			categoryProducts = append(categoryProducts, productInstance)

			b.Request.Visit(categoryLink)
		})
	})

	c.Visit(categoryLink)

	// Workflow: visit --> wait for response --> grab --> visit --> grab products --> Write to output.json
	// c.Visit("")

	// Handle errors
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Alert when done scraping
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished current product scrape", r.Request.URL)
	})

	c.Wait()
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	scrape()

	// e := echo.New()
	// // Serve to echo
	// e.GET("/scrape", func(f echo.Context) error {
	// 	return f.JSON(http.StatusOK, categoriesList)
	// })

	// e.Logger.Fatal(e.Start(":8000"))
}
