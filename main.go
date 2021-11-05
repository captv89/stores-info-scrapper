package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gocolly/colly/v2"
)

type item struct {
	Category    string
	CategoryURL string
	Item        string
	ItemURL     string
	SKU         string
	ImageURL    string
	Description string
}

func main() {
	fmt.Println("Hello, World!")
	allowedDomain := "technicalshipsupplies.com"
	startUrl := "https://www.technicalshipsupplies.com/categories?sort=alphaasc"
	fmt.Println("Starting to Crawl from:", startUrl)
	crawl(startUrl, allowedDomain)
}

func crawl(startUrl string, allowedDomain string) {
	itemArray := []item{}
	c := colly.NewCollector()
	var category, catLink string
	// var prodLink string
	c.IgnoreRobotsTxt = true

	c.SetRequestTimeout(20 * time.Second)

	c.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 3 * time.Second,
		//Parallelism
		//Parallelism: 2,
	})

	// Copy the Collector to a new variable
	d := c.Clone()
	g := c.Clone()

	c.OnHTML("div.SubCategoryList", func(e *colly.HTMLElement) {
		// fmt.Println(e.Text)
		e.ForEach("li > a", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			category = elem.Text
			catLink = link
			// fmt.Printf("Category %s Link: %s", temp.Category, link)
			d.Visit(link)
			// itemStore = append(itemStore, temp)
		})
	})

	d.OnHTML("#LayoutColumn2", func(e *colly.HTMLElement) {

		e.ForEach("div.ProductDetails > strong > a", func(_ int, elem *colly.HTMLElement) {
			prodLink := elem.Attr("href")
			fmt.Println("ProdLink", prodLink)
			// fmt.Println(elem)
			// fmt.Println(prodLink)
			g.Visit(prodLink)
		})
		// fmt.Println("Item Link:", temp.Item, "  ", prodLink)
	})

	g.OnHTML("#ProductDetails > div.BlockContent", func(e *colly.HTMLElement) {
		temp := item{}
		temp.Category = category
		temp.CategoryURL = catLink
		temp.Item = e.ChildText("h2")
		temp.ItemURL = e.Request.URL.String()
		temp.SKU = e.ChildText("span.VariationProductSKU")
		temp.ImageURL = e.ChildAttr("img", "src")
		temp.Description = e.ChildAttr("img", "alt")
		fmt.Println(temp)
		itemArray = append(itemArray, temp)
		writeToFile(itemArray)
	})

	g.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	d.OnHTML("div.FloatRight", func(e *colly.HTMLElement) {
		nextPage := e.ChildAttr("a", "href")
		if nextPage != "" {
			fmt.Println("Next Page:", nextPage)
			d.Visit(nextPage)
		}
	})
	// Commented as there is no pagination for the main list of categories.
	// c.OnHTML("div.FloatRight", func(e *colly.HTMLElement) {
	// 	nextPage := e.ChildAttr("a", "href")
	// 	fmt.Println("Next Category Page Link: ", nextPage)
	// 	c.Visit(nextPage)
	// })

	// Start scraping on url provided
	c.Visit(startUrl)

	// Wait until threads are finished
	c.Wait()
}

func writeToFile(itemData []item) {
	jsonData, err := json.Marshal(itemData)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("./spare-data.json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
